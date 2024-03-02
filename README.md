# Chat

Chat is a powerful software that allows you to easily integrate chats for WhatsApp and Telegram, both text and audio, with an artificial intelligence chat bot. It relies on Tibula RDBMS, providing you the ability to add new users and personalize system prompts and translations directly from the Tibula web framework.

## Requirements

To use this software, you need the following:

1. **ffmpeg**: For audio conversion.
2. **Google Credentials**: For Automatic Speech Recognition (ASR) and Text-to-Speech (TTS).
3. **OpenAI Token**: For AI processing.
4. **Telegram Token**: To integrate with Telegram.
5. **Meta Credentials**: To integrate with WhatsApp.

## Features

- Integration with WhatsApp and Telegram.
- Support for both text and audio chats.
- AI chat bot powered by OpenAI.
- User management and system prompt personalization via Tibula RDBMS.
- Audio conversion using ffmpeg.
- ASR and TTS capabilities using Google Cloud Services.

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
