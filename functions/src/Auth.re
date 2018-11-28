module Firebase = {
  type request;
  type response;
  [@bs.send] external send : (response, Js.Json.t) => Js.Promise.t(unit) = "send";

  type void;

  type handler = (request, response) => Js.Promise.t(unit);

  [@bs.val] [@bs.module "firebase-functions"] [@bs.scope "https"] external onRequest : handler => void = "onRequest";
};

let signIn : Firebase.void = Firebase.onRequest((req, res) => {
  "send a message!"
  |> Js.Json.string
  |> Firebase.send(res);
});
