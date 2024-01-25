# Étape de build
FROM golang:1.22rc2 AS builder

# Définis le répertoire de travail
WORKDIR /app

# Copie les fichiers de dépendances et télécharge les dépendances
COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify

# Copie le code source dans l'image
COPY . .

# Compile l'application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Étape de création de l'image finale
FROM scratch

# Copie l'exécutable depuis l'étape de build
COPY --from=builder /app/main .

# Définis les variables d'environnement nécessaires
ENV WEBHOOK_TOKEN=your_webhook_token
ENV FB_TOKEN=your_fb_token
ENV NOCO_URL=your_noco_url
ENV NOCO_TABLE_ID=your_noco_table_id
ENV NOCO_API_KEY=your_noco_api_key
ENV USER1=your_user1
ENV USER2=your_user2

# Expose le port 8080
EXPOSE 8080

# Commande pour lancer l'application
CMD ["./main"]
