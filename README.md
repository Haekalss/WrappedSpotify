# WrappedSpotify

A Spotify Wrapped clone built with Go.

## Local Development

1. Clone the repository
2. Copy `.env.example` to `.env` and fill in your Spotify credentials
3. Run `go mod download` to install dependencies
4. Run `go run main.go` to start the server

## Deployment to Railway

1. Push your code to GitHub
2. Connect your repository to Railway
3. Set the following environment variables in Railway dashboard:
   - `SPOTIFY_CLIENT_ID`: Your Spotify App Client ID
   - `SPOTIFY_CLIENT_SECRET`: Your Spotify App Client Secret
   - `SPOTIFY_REDIRECT_URI`: Your Railway app URL + `/auth/callback`
   - `FRONTEND_URL`: Your frontend URL
   - `PORT`: Railway will set this automatically

4. Railway will automatically detect the Dockerfile and deploy your app

## Required Environment Variables

- `SPOTIFY_CLIENT_ID`: Your Spotify application client ID
- `SPOTIFY_CLIENT_SECRET`: Your Spotify application client secret  
- `SPOTIFY_REDIRECT_URI`: Callback URL for Spotify authentication
- `FRONTEND_URL`: URL of your frontend application
- `PORT`: Port to run the server on (optional, defaults to 8080)