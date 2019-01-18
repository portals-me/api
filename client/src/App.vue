<template>
  <v-app>
    <router-view name="sidebar" />
    <router-view name="topbar" />
    <v-content>
      <router-view />
    </v-content>
  </v-app>
</template>

<script>
export default {
  data () {
    return {
      user: null,
    }
  },
  methods: {
  },
  mounted: async function () {
  }
}
</script>
