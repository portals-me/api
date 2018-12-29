<template>
  <v-toolbar
    absolute
    app
    clipped-left
    dense
    flat
    dark
    class="indigo darken-1"
  >
    <v-toolbar-side-icon
      @click.stop="drawer = !drawer"
    ></v-toolbar-side-icon>
    <v-toolbar-title><router-link to="/" style="color: #fff; text-decoration: none;">Portals@me</router-link></v-toolbar-title>

    <v-spacer></v-spacer>

    <v-toolbar-items v-if="user != null">
      <v-menu offset-y>
        <v-btn
          slot="activator"
          flat
        >
          <v-avatar color="orange" size="32px">
            <v-img :src="user.picture" />
          </v-avatar>
          &nbsp;&nbsp;{{ user.display_name }}
        </v-btn>
        <v-list>
          <v-list-tile
            @click="signOut"
          >
            <v-list-tile-title>Sign Out</v-list-tile-title>
          </v-list-tile>
        </v-list>
      </v-menu>
    </v-toolbar-items>
  </v-toolbar>
</template>

<script>
export default {
  data () {
    return {
      user: null,
    };
  },
  methods: {
    async signOut () {
      localStorage.setItem('id_token', '');
      localStorage.setItem('user', '{}');
      this.user = null;
      this.$router.push('/signin');
    },
  },
  async mounted () {
    this.user = JSON.parse(localStorage.getItem('user'));

    if (!this.user) {
      this.$router.push('/');
    }
  },
}
</script>

