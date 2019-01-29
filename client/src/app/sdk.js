import axios from 'axios';
import api from '../../../src/api';
const url = 'https://v6bnqbi2hf.execute-api.ap-northeast-1.amazonaws.com/prod';
const sdk = api.genSDK(url, () => localStorage.getItem('id_token'), axios);

export default sdk;
