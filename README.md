# 🍳 Recipe Bot

**Recipe Bot** is a Telegram bot that transforms your video content into beautifully formatted recipes using  OpenAI! 🚀 It’s your ultimate kitchen assistant, making recipe creation seamless and fun. 🎉

## ✨ Features

- 🎥 **Video Detection**: Automatically detects video links in your messages.
- 🗣️ **Speech-to-Text Conversion**: Converts video audio into text using OpenAI’s advanced technology.
- 📜 **Recipe Generation**: Creates stunning recipes with titles, bodies, and markdown formatting.
- 📂 **Recipe Storage**: Saves recipes with metadata like cover images for easy access.

## 🛠️ How It Works

1.  The bot listens for your text messages in Telegram.
2.  It checks if the message contains a video link.
3.  If a video is found:
   - Downloads the video’s audio. 
   - Converts the audio to text. 
   - Processes the text to generate a recipe. 
4.  The recipe is sent back to you in a beautifully formatted markdown style.
## 🔑 Prerequisites

- 🐹 **Go**: Version 1.23 or higher.
- 🐳 **Docker**: For building and deploying the application.
- 🤖 **Telegram Bot API**: A bot token to interact with Telegram.
- 🧠 **OpenAI API**: For speech-to-text and recipe generation.
- 📦 **FFmpeg**: For audio processing.
- 🌐 **TikHub API**: For accessing TikTok video data.