<template>
  <form>
    <v-flex>
      <v-text-field
        v-model="form.name"
        label="ユーザーID"
        append-outer-icon="check"
      />
      <v-text-field
        v-model="form.display_name"
        label="表示される名前"
      />
    </v-flex>
    <v-avatar color="orange" size="32px">
      <v-img :src="form.picture" />
    </v-avatar>
    <v-btn depressed>アイコンをアップロード</v-btn>

    <br />

    <v-btn color="success" @click="$emit('submit', form)">送信</v-btn>
    <v-btn depressed @click="$emit('cancel')">キャンセル</v-btn>
  </form>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  props: [
    'formData',
  ],
  data () {
    return {
      form: {
        name: this.formData.name,
        display_name: this.formData.display_name,
        picture: this.formData.picture,
      },
    };
  },
})
</script>

