import axios from 'axios';
const url = 'https://fubo4rnsod.execute-api.ap-northeast-1.amazonaws.com/dev';
const token = localStorage.getItem('id_token');

const sdk = {
  signIn: async (gtoken) => {
    const result = await axios.post(`${url}/signIn`, gtoken);
    return result.data;
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
    list: async () => {
      const result = await axios.get(`${url}/projects`, { headers: { Authorization: `Bearer ${token}` } });
      return result.data;
    },
  },
};

export default sdk;
