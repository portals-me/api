import firebase from 'firebase/app';
import 'firebase/firestore';
import app from '@/instance/firebase';

const firestore = firebase.firestore(app);
const settings = {
  timestampsInSnapshots: true,
};
firestore.settings(settings);

export default firestore;
