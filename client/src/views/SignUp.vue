<template>
  <v-container>
    <v-layout
      text-xs-center
      wrap
    >
      <v-flex mb-5 xs6>
        <v-form>
          <v-text-field
            v-model="name"
            label="User ID"
            required
          />

          <v-text-field
            v-model="display_name"
            label="Display Name"
            required
          />

          <v-avatar color="orange" size="32px">
            <v-img :src="photo" />
          </v-avatar>

          <v-btn>
            Upload Icon
          </v-btn>
        </v-form>
        <p>{{ errorMessage }}</p>
        <v-btn @click="signUp">Create an account</v-btn>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import sdk from '@/app/sdk';

export default {
  data () {
    return {
      name: '',
      display_name: '',
      photo: '',
      errorMessage: '',
    };
  },
  methods: {
    async signUp () {
      const user = await this.$gAuth.signIn();

      try {
        const result = (await sdk.signUp({
          google_token: user.getAuthResponse().id_token,
          name: this.name,
          display_name: this.display_name,
          photo: this.photo,
        })).data;

        localStorage.setItem('id_token', result.id_token);
        localStorage.setItem('user', JSON.stringify(result.user));
        this.$router.push('/');
      } catch (err) {
        this.errorMessage = err.response.data;
        return;
      }
    },
    async onMounted () {
      this.name = this.$route.query.name;
      this.display_name = this.$route.query.display_name;
      this.photo = this.$route.query.photo;
    },
  },
  async mounted () {
    await this.onMounted();
  },
}
</script>
