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
      <v-toolbar-title>Portals@me</v-toolbar-title>

      <v-spacer></v-spacer>

      <v-toolbar-items>
        <v-btn flat>
          <v-avatar color="orange" size="32px">
            <span class="white--text headline">A</span>
          </v-avatar>
          &nbsp;&nbsp;{{ user.user_name }}
        </v-btn>
      </v-toolbar-items>
    </v-toolbar>

    <v-content>
      <router-view />
    </v-content>
  </v-app>
</template>

<script>
import HelloWorld from './components/HelloWorld'
import * as user from '@/sdk/user';

export default {
  name: 'App',
  components: {
    HelloWorld
  },
  data () {
    return {
      drawer: null,
      user: {},
    }
  },
  methods: {
    async loadUser () {
      this.user = await user.api.getMe();
    }
  },
  mounted: async function () {
    await this.loadUser();
  }
}
</script>
