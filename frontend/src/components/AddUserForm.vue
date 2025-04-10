<template>
    <v-container>
      <v-form @submit.prevent="createUser">
        <v-text-field
          label="Username"
          v-model="username"
          required
        ></v-text-field>
        <v-btn type="submit" color="primary">Add User</v-btn>
      </v-form>
      <v-alert v-if="successMessage" type="success" class="mt-4">
        {{ successMessage }}
      </v-alert>
    </v-container>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue'
  
  // Define types
  interface User {
    id: number
    username: string
  }
  
  const username = ref<string>('')
  const successMessage = ref<string>('')
  
  const createUser = async (): Promise<void> => {
    try {
      const res = await fetch('http://localhost:8080/users', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username: username.value }),
      })
  
      if (res.ok) {
        const user: User = await res.json()
        successMessage.value = `User ${user.username} created with ID ${user.id}`
        username.value = ''
      } else {
        const errorText = await res.text()
        console.error('Create failed:', errorText)
        alert('Failed to create user')
      }
    } catch (error) {
      console.error('Request error:', error)
      alert('An error occurred')
    }
  }
  </script>
  