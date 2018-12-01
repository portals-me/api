<template>
  <v-card>
    <v-img
      :src="$props.ogp.image"
      max-height="200px"
    />

    <v-card-title primary-title>
      <div>
        <h3 class="headline mb-0">{{ $props.ogp.title }}</h3>
        <div>{{ $props.ogp.description }}</div>
      </div>
    </v-card-title>

    <v-card-actions>
      <v-btn :href="$props.ogp.url" flat color="orange">Open</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: 'ogp-card',
  props: [ 'ogp' ],
}
</script>

