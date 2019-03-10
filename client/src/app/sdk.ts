import axios from 'axios';
const url = process.env.VUE_APP_API_ENDPOINT;
const getToken = () => localStorage.getItem('id_token');

export default {
  signUp: async (json: object) => {
    return await axios.post(`${url}/auth/signUp`, JSON.stringify(json));
  },
  signIn: async (json: object) => {
    return await axios.post(`${url}/auth/signIn`, JSON.stringify(json));
  },
  user: {
    me: async () => {
      return await axios.get(`${url}/users/me`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    get: async (userName: string) => {
      return await axios.get(`${url}/users/${userName}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    update: async (userId: string, form: any) => {
      return await axios.put(`${url}/users/${userId}`, form, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    follow: async (userName: string) => {
      return await axios.post(`${url}/users/${userName}/follow`, {}, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    feed: {
      list: async (userName: string) => {
        return await axios.get(`${url}/users/${userName}/feed`, { headers: { Authorization: `Bearer ${getToken()}` } });
      },
    },
  },
  comment: {
    create: async (collectionId: string, message: object) => {
      return await axios.post(`${url}/comments`, {
        collectionId,
        message,
      }, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    list: async (collectionId: string) => {
      return await axios.get(`${url}/collections/${collectionId}/comments`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
  },
  timeline: {
    get: async () => {
      return await axios.get(`${url}/timeline`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
  },
  collection: {
    create: async (form: object) => {
      return await axios.post(
        `${url}/collections`,
        form,
        { headers: { Authorization: `Bearer ${getToken()}` } }
      );
    },
    get: async (collectionId: string) => {
      return await axios.get(`${url}/collections/${collectionId}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    list: async () => {
      return await axios.get(`${url}/collections`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    delete: async (collectionId: string) => {
      return await axios.delete(`${url}/collections/${collectionId}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
  },
  article: {
    create: async (collectionId: string, form: object) => {
      return await axios.post(
        `${url}/collections/${collectionId}/articles`,
        form,
        { headers: { Authorization: `Bearer ${getToken()}` } }
      );
    },
    generate_presigned_url: async (collectionId: string, key: object) => {
      return await axios.post(
        `${url}/collections/${collectionId}/articles-presigned`,
        key,
        { headers: { Authorization: `Bearer ${getToken()}` } }
      );
    },
    list: async (collectionId: string) => {
      return await axios.get(`${url}/collections/${collectionId}/articles`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
  },
};
