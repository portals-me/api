<template>
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
</template>

<script>
export default {
  data () {
    return {
      drawer: true,
    };
  }
}
</script>

