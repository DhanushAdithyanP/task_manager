<template>
  <v-container class="fill-height d-flex justify-center align-center">
    <v-card class="pa-5" width="500">
      <v-card-title class="text-h6">Login</v-card-title>
      <v-card-text>
        <v-text-field v-model="username" label="Username" outlined />
        <v-text-field v-model="password" label="Password" type="password" outlined />
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" @click="login">Login</v-btn>
      </v-card-actions>
      <v-alert v-if="error" type="error" dense>{{ error }}</v-alert>
    </v-card>
  </v-container>
</template>

  
  <script setup>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  
  const username = ref('')
  const password = ref('')
  const error = ref('')
  const router = useRouter()
  
  const login = async () => {
  try {
    const res = await fetch('http://localhost:8080/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value, password: password.value }),
      credentials: 'include' 
    })

    if (!res.ok) throw new Error('Invalid login')

    const data = await res.json()
    localStorage.setItem('token', data.token)
    localStorage.setItem('role', data.role)
    localStorage.setItem('userId', data.user_id) 

    if (data.role === 'admin') {
      router.push('/admin')
    } else {
      router.push('/dashboard')
    }
  } catch (err) {
    error.value = err.message
  }
}

  </script>
  