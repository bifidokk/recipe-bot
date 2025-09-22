# 🍳 Recipe Bot

**Recipe Bot** is a Telegram bot that transforms your video content into beautifully formatted recipes using  OpenAI! 🚀 It’s your ultimate kitchen assistant, making recipe creation seamless and fun. 🎉

## ✨ Features

- 🎥 **Video Detection**: Automatically detects video links in your messages.
- 🎬 **Multi-Platform Support**: Works with TikTok videos and Instagram Reels/Posts.
- 🗣️ **Speech-to-Text Conversion**: Converts video audio into text using OpenAI.
- 📜 **Recipe Generation**: Creates stunning recipes with titles, bodies, and markdown formatting.
- 📂 **Recipe Storage**: Saves recipes with metadata for easy access.

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
- 🐍 **Python 3**: For Instagram support via Instaloader.
- 📱 **Instaloader**: For downloading Instagram content.

## 📱 Supported Platforms

### TikTok
- ✅ Short URLs: `https://vm.tiktok.com/ABC123`
- ✅ Full URLs: `https://www.tiktok.com/@user/video/1234567890`
- ✅ Uses TikHub API for video data extraction

### Instagram
- ✅ Reels: `https://www.instagram.com/reel/ABC123/`
- ✅ Posts: `https://www.instagram.com/p/XYZ789/`
- ✅ Uses Instaloader for content download
- ✅ Requires Python 3 and Instaloader installation