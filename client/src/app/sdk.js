import axios from 'axios';
const url = 'https://fubo4rnsod.execute-api.ap-northeast-1.amazonaws.com/dev';

const sdk = {
  signIn: async (gtoken) => {
    const result = await axios.post(`${url}/signIn`, gtoken);
    return result.data;
  },
};

export default sdk;
