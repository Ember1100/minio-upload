<script setup>
import { UploadFilled } from '@element-plus/icons-vue'

import md5 from "../lib/md5";
import { taskInfo, initTask, preSignUrl, merge } from '../lib/api';
import {ElNotification} from "element-plus";
import Queue from 'promise-queue-plus';
import axios from 'axios'
import { ref } from 'vue'

// 文件上传分块任务的队列（用于移除文件时，停止该文件的上传队列） key：fileUid value： queue object
const fileUploadChunkQueue = ref({}).value

/**
 * 获取一个上传任务，没有则初始化一个
 */
const getTaskInfo = async (file) => {
    let task;
    const identifier = await md5(file)
    const { code, data, msg } = await taskInfo(identifier)
    if (code === 200000) {
        task = data
        if (!task) {
            const initTaskData = {
                identifier,
                fileName: file.name,
                totalSize: file.size,
                chunkSize: 5 * 1024 * 1024
            }
            const { code, data, msg } = await initTask(initTaskData)
            if (code === 200000) {
                task = data
            } else {
                ElNotification.error({
                    title: '文件上传错误',
                    message: msg
                })
            }
        }
    } else {
        ElNotification.error({
            title: '文件上传错误',
            message: msg
        })
    }
    return task
}

/**
 * 上传逻辑处理，如果文件已经上传完成（完成分块合并操作），则不会进入到此方法中
 */
const handleUpload = (file, taskRecord, options) => {

    let uploadedSize = 0 // 已上传的大小
    const totalSize = file.size || 0 // 文件总大小
    const { exitPartList, chunkSize, chunkNum, fileIdentifier } = taskRecord

    const uploadNext = async (partNumber) => {
        const start = new Number(chunkSize) * (partNumber - 1)
        const end = start + new Number(chunkSize)
        const blob = file.slice(start, end)
        const { code, data, msg } = await preSignUrl({ identifier: fileIdentifier, partNumber: partNumber} )
        if (code === 200000 && data) {
            await axios.request({
                url: data,
                method: 'PUT',
                data: blob,
                headers: {
                    'Content-Type': 'application/octet-stream'
                }
            })
            return Promise.resolve({ partNumber: partNumber, uploadedSize: blob.size })
        }
        return Promise.reject(`分片${partNumber}， 获取上传地址失败`)
    }

    /**
     * 更新上传进度
     * @param increment 为已上传的进度增加的字节量
     */
    const updateProcess = (increment) => {
        increment = new Number(increment)
        const { onProgress } = options
        let factor = 1000; // 每次增加1000 byte
        let from = 0;
        // 通过循环一点一点的增加进度
        while (from <= increment) {
            from += factor
            uploadedSize += factor
            const percent = Math.round(uploadedSize / totalSize * 100).toFixed(2);
            onProgress({percent: percent})
        }
    }


    return new Promise(resolve => {
        const failArr = [];
        const queue = Queue(5, {
            "retry": 3,               //Number of retries
            "retryIsJump": false,     //retry now?
            "workReject": function(reason,queue){
                failArr.push(reason)
            },
            "queueEnd": function(queue){
                resolve(failArr);
            }
        })
        fileUploadChunkQueue[file.uid] = queue
        for (let partNumber = 1; partNumber <= chunkNum; partNumber++) {
            const exitPart = (exitPartList || []).find(exitPart => exitPart.partNumber == partNumber)
            if (exitPart) {
                // 分片已上传完成，累计到上传完成的总额中
                updateProcess(exitPart.size)
            } else {
                queue.push(() => uploadNext(partNumber).then(res => {
                    // 单片文件上传完成再更新上传进度
                    updateProcess(res.uploadedSize)
                }))
            }
        }
        if (queue.getLength() == 0) {
            // 所有分片都上传完，但未合并，直接return出去，进行合并操作
            resolve(failArr);
            return;
        }
        queue.start()
    })
}

/**
 * el-upload 自定义上传方法入口
 */
const handleHttpRequest = async (options) => {
    const file = options.file
    const task = await getTaskInfo(file)
    if (task) {
        const { finished, path, taskRecord } = task
        const { fileIdentifier: identifier } = taskRecord
        if (finished) {
            return path
        } else {
            const errorList = await handleUpload(file, taskRecord, options)
            if (errorList.length > 0) {
                ElNotification.error({
                    title: '文件上传错误',
                    message: '部分分片上次失败，请尝试重新上传文件'
                })
                return;
            }
            const { code, data, msg } = await merge(identifier)
            if (code === 200000) {
                return path;
            } else {
                ElNotification.error({
                    title: '文件上传错误',
                    message: msg
                })
            }
        }
    } else {
        ElNotification.error({
            title: '文件上传错误',
            message: '获取上传任务失败'
        })
    }
}

/**
 * 移除文件列表中的文件
 * 如果文件存在上传队列任务对象，则停止该队列的任务
 */
const handleRemoveFile = (uploadFile, uploadFiles) => {
    const queueObject = fileUploadChunkQueue[uploadFile.uid]
    if (queueObject) {
        queueObject.stop()
        fileUploadChunkQueue[uploadFile.uid] = undefined
    }
}

</script>
<template>
    <el-card style="width: 80%; margin: 80px auto" header="文件分片上传">
        <el-upload
            class="upload-demo"
            drag
            action="/"
            multiple
            :http-request="handleHttpRequest" 
            :on-remove="handleRemoveFile">
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
                请拖拽文件到此处或 <em>点击此处上传</em>
            </div>
        </el-upload>
    </el-card>

</template>
