# Garconia Law Bot

Welcome to the official Garconia Law Bot

---

## What is Garconia?

Garconia is a **Minecraft server** inspired by modern-day democracy. In Garconia, you can:

- Start your own government.
- Work with the local government.
- Become the opposition party.

With less emphasis on PvP, our server is perfect for players who enjoy a wide range of gameplay styles, whether you're a builder, a strategist, or a diplomat.

---

## How Do I Join?

Our **Discord server** is the best place to learn more about Garconia, connect with the community, and get updates.

Join us here: [Discord Invite](https://discord.gg/T7nbWurxcT)

---

## What is This Repository?

This repository contains the code for the **Garconia Law Bot**. The Garconia Law Bot handles showcasing the users the intricate constitution of Garconia, with multiple commands to:

1. `Get All Articles`
2. `Get Specific Articles`
3. `Get Specific Clauses`
4. `Get Specific Amendments`
5. And More...

---

## Get Started

To get started with the Garconia Law Bot:

1. Clone the repository:
   ```bash
   git clone https://github.com/Arinji2/garconia-law-bot
   ```

2. Set up your environment:

   - Rename the `example.env` file to `.env`.
   - Populate the `.env` file with the following required variables:

     ```env
     TOKEN=(Bot Token)
     GUILD_ID=(ID of the Guild which the bot is in)
     ADMIN_EMAIL=(Super user email of the Pocketbase Instance)
     ADMIN_PASSWORD=(Super user password of the Pocketbase Instance)
     BASE_DOMAIN=(Base domain of the Pocketbase Instance)
     ALLOWED_ROLES=(Comma-separated list of allowed role IDs for the data-refresh command)
     ALLOWED_CHANNELS=(Comma-separated list of allowed channel IDs for the commands to be run in)
     ```

3. Install dependencies and run the bot:
   ```bash
   go mod tidy
   go run .
   ```

---

## Commands

The following commands exist in the bot:

### Refresh Data

- **Command:** `/refresh-data`
- **Description:** Refresh the data of the bot.

### Get Clauses

- **Command:** `/get-clauses`
- **Description:** Get the clauses of the constitution.
- **Options:**
  - `article-number` (Required): Article number of the clause.
  - `clause-number` (Optional): Clause number (autocomplete enabled).

### Get Articles

- **Command:** `/get-articles`
- **Description:** Get the articles of the constitution.
- **Options:**
  - `article-number` (Optional): Article number of the clause (autocomplete enabled).

### Get Amendments

- **Command:** `/get-amendments`
- **Description:** Get the amendments of the constitution.
- **Options:**
  - `article-number` (Required): Article number of the constitution.
  - `clause-number` (Required): Clause number of the article.
  - `amendment-number` (Required): Amendment number of the clause.

---

## Built by Arinji

This project was proudly built by [Arinji](https://www.arinji.com/). Check out my website for more cool projects!

