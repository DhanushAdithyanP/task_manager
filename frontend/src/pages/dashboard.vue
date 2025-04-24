<template>
    <v-container>
      <v-row>
        <v-col cols="12">
          <h2 class="text-h5 mb-4">Welcome to Your Dashboard</h2>
  
          <!-- Top 5 Tasks Table -->
          <v-card class="mb-4">
            <v-card-title>Top 5 Priority Tasks</v-card-title>
            <v-card-text>
              <v-data-table-virtual
                :headers="topTaskHeaders"
                :items="topTasks"
                item-value="title"
                height="300"
                fixed-header
                dense
              />
            </v-card-text>
          </v-card>
  
          <!-- Create New Task -->
          <v-card class="mb-4">
            <v-card-title>Add New Task</v-card-title>
            <v-card-text>
              <v-text-field v-model="newTask.title" label="Title" />
              <v-text-field v-model.number="newTask.priority" label="Priority" type="number" />
              <v-btn color="primary" @click="addTask">Add Task</v-btn>
            </v-card-text>
          </v-card>
  
          <!-- All Tasks Table -->
          <v-card>
            <v-card-title>Your Tasks</v-card-title>
            <v-card-text>
              <v-data-table
                :headers="taskHeaders"
                :items="tasks"
                item-value="id"
                show-expand
                class="elevation-1"
                dense
              >
                <!-- Expandable Row: Subtasks -->
                <template v-slot:expanded-row="{ item }">
                  <v-list>
                    <v-list-item v-for="subtask in subtasks[item.id]" :key="subtask.task_id">
                      <v-list-item-content>
                        <v-list-item-title>{{ subtask.title }}</v-list-item-title>
                      </v-list-item-content>
                    </v-list-item>
                  </v-list>
                </template>
  
                <!-- Expand Button -->
                <template v-slot:expanded-item="{ item }">
                  <v-btn @click="fetchSubtasks(item.id)">Load Subtasks</v-btn>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import axios from 'axios'
  
  const userId = localStorage.getItem('userId')
  const username = localStorage.getItem('username') || 'User'
  
  const tasks = ref([])
  const subtasks = ref({})
  const topTasks = ref([])
  
  const newTask = ref({
    title: '',
    priority: 1
  })
  
  const taskHeaders = [
    { text: 'Task', value: 'title' },
    { text: 'Priority', value: 'priority' }
  ]
  
  const topTaskHeaders = [
    { title: 'Task', key: 'title' },
    { title: 'Priority', key: 'priority' }
  ]
  
  const fetchTasks = async () => {
    const response = await axios.get(`http://localhost:8080/tasks/${userId}`, {
      withCredentials: true
    })
    tasks.value = response.data
  }
  
  const fetchSubtasks = async (taskId) => {
    if (!subtasks.value[taskId]) {
      const res = await axios.get(`http://localhost:8080/tasks/${taskId}/subtasks`, {
        withCredentials: true
      })
      subtasks.value[taskId] = res.data
    }
  }
  
  const fetchTopTasks = async () => {
    const res = await axios.get(`http://localhost:8080/users/${userId}/top-tasks`, {
      withCredentials: true
    })
    topTasks.value = res.data
  }
  
  const addTask = async () => {
    const taskPayload = { ...newTask.value, user_id: parseInt(userId) }
    await axios.post('http://localhost:8080/tasks', taskPayload, {
      withCredentials: true
    })
    newTask.value.title = ''
    newTask.value.priority = 1
    await fetchTasks()
    await fetchTopTasks()
  }
  
  onMounted(() => {
    fetchTasks()
    fetchTopTasks()
  })
  </script>
  
  <style scoped>
  .v-card {
    margin-bottom: 20px;
  }
  
  .text-subtitle-1 {
    font-size: 1.25rem;
  }
  </style>
  