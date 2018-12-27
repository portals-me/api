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
      drawer: null,
      user: null,
    }
  },
  methods: {
    async signOut () {
      localStorage.setItem('id_token', '');
      localStorage.setItem('user', '{}');
      this.user = null;
      this.$router.push('/signin');
    },
    async onMount () {
      this.user = JSON.parse(localStorage.getItem('user'));

      if (!this.user) {
        this.$router.push('/signin');
      }
    },
  },
  mounted: async function () {
    await this.onMount();
  }
}
</script>
