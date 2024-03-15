# Chat

Chat is a powerful software that allows you to easily integrate text and audio chats for WhatsApp, Telegram and Asterisk, with an artificial intelligence chat bot. It relies on Tibula RDBMS, providing you the ability to add new users and personalize system prompts and translations directly from the Tibula web framework.

## Requirements

To use this software, you need the following:

- **FFmpeg**: For audio conversion.
- **OpenAI Token**: For AI and speech processing.
- **Google Credentials**: For Automatic Speech Recognition (ASR) and Text-to-Speech (TTS).
- **Telegram Token**: To integrate with Telegram.
- **Meta Credentials**: To integrate with WhatsApp.
- **Asterisk Token**: To integrate with Asterisk.

## Features

- Integration with WhatsApp, Telegram and Asterisk.
- Support for both text and audio chats.
- AI chat bot powered by OpenAI.
- User management and system prompt personalization via Tibula RDBMS.
- Audio conversion using ffmpeg.
- ASR and TTS capabilities using OpenAI (default) or Google Cloud Services.

## Installation

```
git clone https://github.com/eja/chat
cd chat
make
./chat --wizard
./chat --start
```

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
