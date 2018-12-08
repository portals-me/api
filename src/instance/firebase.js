import firebase from 'firebase';

const app = firebase.initializeApp({
  apiKey: "AIzaSyAL-NuKxMhShZVARoxzNMvXrGN3A65OEps",
  authDomain: "portals-me.firebaseapp.com",
  databaseURL: "https://portals-me.firebaseio.com",
  projectId: "portals-me",
  storageBucket: "portals-me.appspot.com",
  messagingSenderId: "670077302427"
});

export default app;
