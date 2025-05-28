hmmmm 

ok having a django page for deployments:
- I want to know
  - which app server is handling my session (the version of the software and the running process/vm/host)
  - how to set affinity for a particular app server (e.g. show me a list of all running servers and let me connect to one directly and consistently).
    - e.g. /home/settings?affinity=app-v1.23-i1
    - e.g. /home/settings?affinity=app-v1.23-i3
  - status of the deployment
  - history of deployments

What is a secure way to distribute VPN details? Is there a way to do it through
a static site I update asynchronously?
- for when I want to administer/configure wireguard for a team.
- vpn.couetil.com 
  -> instructions to configure wireguard
  -> the vpn server needs the user's public key.
  -> the user needs the vpn's public key.
  -> the vpn server has to load the new configuration with the client's pubkey
  added
  -> the user needs to update their configuration with the server's pubkey
- I could gate vpn.couetil.com with auth.couetil.com and use some kind of SSO
  in combination with permissions based on SSO/auth role.
  -> then client will POST public key to file, which will respond with a 
    -> GET vpn.couetil.com/
    <- 200 Ok
    -> POST vpn.couetil.com/clientpubkey
    <- 303 See Other, Location: /pubkey
    -> GET vpn.couetil.com/pubkey
    <- 200 Content-Disposition: attachment; Content-Type: text/plain
  -> then client can download a file that clicking on will load into wireguard?


What does an Authentication/Authorization flow look like?
-> [client] GET admin.couetil.com/
<- [admin] 307 Temporary Redirect, Location: authn.couetil.com/oauth2/authorize
-> [client] GET authn.couetil.com/oauth2/authorize
      ?response_type=code
      &client_id=CLIENT_ID
      &redirect_uri=https://admin.couetil.com/authn/callback
      &scope=openid+profile+email
      &state=abc123
<- [authn] 200 Ok
-> [client] POST authn.couetil.com/credentials
<- [authn] 303 See Other,
      Location: admin.couetil.com/authn/callback?code=CODE&state=
-> [client] GET admin.couetil.com/authn/callback?code=CODE&state=
  -> [admin] POST https://authn.couetil.com/oauth2/token
        ?grant_type=authorization_code
        &code=CODE
        &redirect_uri=https://admin.example.com/callback
        Authorization: Basic <preset user/password for admin-authn communication>
  <- [authn] 200 Ok, { access_token, id_token, token_type, expires_in }
<- [admin] 307 Temporary Redirect,
      Location: /
      Set-Cookie: <session_authorization_token>


    -> GET vpn.couetil.com/
    <- 307 Temporary Redirect, Location: auth.couetil.com
    -> GET auth.couetil.com/
