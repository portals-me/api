import axios from 'axios';
const url = 'https://1cmunxpbgi.execute-api.ap-northeast-1.amazonaws.com/dev';
const token = localStorage.getItem('id_token');

const sdk = {
  signIn: async (gtoken) => {
    const result = await axios.post(`${url}/signIn`, gtoken);
    return result.data;
  },
  user: {
    me: async () => {
      const result = await axios.get(`${url}/users/me`, { headers: { Authorization: `Bearer ${token}` } });
      return result.data;
    },
  },
  comment: {
    create: async (projectId, message) => {
      const result = await axios.post(`${url}/comments`, {
        projectId,
        message,
      }, { headers: { Authorization: `Bearer ${token}` } });
      return result.data;
    },
    list: async (projectId) => {
      const result = await axios.post(`${url}/projects/${projectId}/comments`, {
        projectId,
      }, { headers: { Authorization: `Bearer ${token}` } });
      return result.data;
    },
  },
  project: {
    create: async (form) => {
      await axios.post(
        `${url}/projects`,
        form,
        { headers: { Authorization: `Bearer ${token}` } }
      );

      return;
    },
    get: async (projectId) => {
      const result = await axios.get(`${url}/projects/${projectId}`, { headers: { Authorization: `Bearer ${token}` } });
      return result.data;
    },
    list: async () => {
      const result = await axios.get(`${url}/projects`, { headers: { Authorization: `Bearer ${token}` } });
      return result.data;
    },
  },
};

export default sdk;
