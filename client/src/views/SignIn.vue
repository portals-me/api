<template>
  <v-container>
    <v-layout
      text-xs-center
      wrap
    >
      <v-flex mb-5 xs12>
        <p style="color: red;">{{ this.errorMessage }}</p>
        <v-btn @click="signIn">Sign In with Google</v-btn>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import sdk from '@/app/sdk';

export default {
  data () {
    return {
      errorMessage: '',
    };
  },
  methods: {
    async signIn () {
      const user = await this.$gAuth.signIn();

      try {
        const result = (await sdk.signIn(user.getAuthResponse().id_token)).data;

        localStorage.setItem('id_token', result.id_token);
        this.$store.commit('setUser', result.user);
        this.$router.push('/');
      } catch (err) {
        this.errorMessage = err.response.data;
        return;
      }
    },
  },
  async mounted () {
  },
}
</script>

