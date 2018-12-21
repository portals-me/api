<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      fixed
      clipped
      app
    >
      <v-list dense class="pt-0">
        <v-subheader>Your Channels</v-subheader>

        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>dashboard</v-icon>
          </v-list-tile-action>

          <v-list-tile-content>
            <v-list-tile-title>Timeline</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>notifications</v-icon>
          </v-list-tile-action>

          <v-list-tile-content>
            <v-list-tile-title>@me</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-subheader>Project Channels</v-subheader>

        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>movie</v-icon>
          </v-list-tile-action>

          <v-list-tile-content>
            <v-list-tile-title>@meow</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>brush</v-icon>
          </v-list-tile-action>

          <v-list-tile-content>
            <v-list-tile-title>@cat</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>more_vert</v-icon>
          </v-list-tile-action>

          <v-list-tile-content>
            <v-list-tile-title>...and more</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-subheader>Shared Channels</v-subheader>

        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>public</v-icon>
          </v-list-tile-action>

          <v-list-tile-content>
            <v-list-tile-title>All</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>comment</v-icon>
          </v-list-tile-action>

          <v-list-tile-content>
            <v-list-tile-title>#programming</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>comment</v-icon>
          </v-list-tile-action>

          <v-list-tile-content>
            <v-list-tile-title>#haskell</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>

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

      <v-toolbar-items v-if="$store.state.user != null">
        <v-menu offset-y>
          <v-btn
            slot="activator"
            flat
          >
            <v-avatar color="orange" size="32px">
              <v-img :src="$store.state.user.picture" />
            </v-avatar>
            &nbsp;&nbsp;{{ $store.state.user.display_name }}
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
    }
  },
  methods: {
    async signOut () {
      localStorage.setItem('id_token', '');
      this.$store.dispatch('signOut');
      this.$router.push('/signin');
    },
    async onMount () {
      await this.$store.dispatch('initialize');
    },
  },
  mounted: async function () {
    await this.onMount();
  }
}
</script>
