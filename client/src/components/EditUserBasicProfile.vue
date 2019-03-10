<template>
  <form>
    <v-flex>
      <v-text-field
        v-model="form.name"
        label="ユーザーID"
        :append-outer-icon="user_id_icon"
        @input="checkUserExists"
      />
      <v-text-field
        v-model="form.display_name"
        label="表示される名前"
      />
    </v-flex>
    <v-avatar color="orange" size="32px">
      <v-img :src="form.picture" />
    </v-avatar>
    <v-btn
      depressed
      @click="uploadIconPicture"
    >アイコンをアップロード</v-btn>

    <input
      type="file"
      accept="image/*"
      @change="selectFiles"
    />
    <div class="preview" v-if="imageData">
      <img
        :src="imageData"
      >
    </div>

    <br />

    <v-btn color="success" @click="$emit('submit', form)">送信</v-btn>
    <v-btn depressed @click="$emit('cancel')">キャンセル</v-btn>
  </form>
</template>

<script lang="ts">
import Vue from 'vue'
import axios from 'axios'
import sdk from '@/app/sdk'

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
      user_id_icon: 'check',
      imageData: null,
    };
  },
  methods: {
    async checkUserExists() {
      if (!this.form.name) {
        this.user_id_icon = 'cancel';
        return;
      }

      await sdk.user.get(this.form.name)
        .then(_ => {
          this.user_id_icon = 'cancel';
        })
        .catch(_ => {
          this.user_id_icon = 'check';
        });
    },
    async selectFiles(input: any) {
      const url = URL.createObjectURL(input.target.files[0]);
      this.imageData = url;
    },
    async uploadIconPicture() {
      const presignedURL = (await sdk.article.generate_presigned_url('user', this.imageData.file.name)).data;
      await axios.put(presignedURL, this.imageData.file, {
        headers: { 'Content-Type': this.imageData.file.type },
      });

      const user = JSON.parse(localStorage.getItem('user') as any);
      this.form.picture = `https://s3-ap-northeast-1.amazonaws.com/portals-me-storage-users/${encodeURIComponent(user.id)}/user/${this.imageData.file.name}`;
    },
  },
})
</script>

