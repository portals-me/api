import axios from 'axios';

export const api = {
  list: async () => {
    return (await axios.get('http://localhost:5000/projects')).data;
  },
  get: async (projectId) => {
    return (await axios.get(`http://localhost:5000/projects/${projectId}`)).data;
  },
};
