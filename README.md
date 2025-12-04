# Artificial Suspects

**Artificial Suspects** is an experimental art game that explores the biases and prejudices of different AI models.
Rather than analyzing training datasets or statistical outputs, this game allows players to experience AI biases firsthand through gameplay in *Unusual Suspects*.

Play the *Unusual Suspects* board game with AI in your browser!

## Contributing

You can contribute to the project by reporting issues, suggesting features, giving feedback and also sharing it on the net!

For code contributions, please check the DEVELOPING chapter.
Project is not fully prepared for comfortable external contributions, but if you are brave enough, feel free to fork and try it out!

## DEVELOPING

### Architecture

Frontend is written in SvelteKit and backend is written in Go.
- only backend calls the LLM services, secrets lies there
- frontend communicates only with backend, no secret keys stored in frontend
- AI model is selected at the start of the game and is used for the whole game

### Frontend server

```
cd front
npm run dev
```

### Backend server
```
cd backend
go run main.go
```

## Deployment

### Build Backend Docker Image

```bash
docker build -t agajdosi/artsus_server:latest --platform linux/amd64  .
docker push agajdosi/artsus_server:latest
```

## Acknowledgments
A huge thanks to **SvelteKit** for making this possible! 🎉
