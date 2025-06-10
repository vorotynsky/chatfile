# Chatfile

**Chatfile** is a simple tool for interacting with language models using plaintext. 

## Example

```
FROM chatgpt
SYSTEM You are a technical writer.

ASK Write a readme file for the chatfile project
ANSWER |
    **Chatfile** is a prompt-driven tool for interacting with LLMs using plaintext. 
```

## Usage

Set your API key and optionally the base url of an openai-compatible api:

```shell
export OPENAI_API_KEY="sk-chatfile-key"
export OPENAI_BASE_URL="http://api.example.com/" # Optional
```

Create a chatfile:

```shell
cat > ./chatfile << ---
FROM gpt-4.1-nano
SYSTEM |
    You are an expert in explaining mathematical concepts clearly and concisely. 
    Provide a straightforward and accurate explanation for students.

ASK What is Cauchy's functional equation?
---
```

Call the tool to execute your chatfile and get a response with this command:

```shell
chatfile run ./chatfile
```

---

**Chatfile** — prompt and get responses all in one file!
