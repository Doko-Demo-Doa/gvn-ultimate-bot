import { Client, GatewayIntentBits } from "discord.js";
const client = new Client({ intents: [GatewayIntentBits.Guilds] });

export async function GET(request: Request) {
  const oauthGuild = await client.guilds.fetch(
    process.env.DISCORD_GUILD_ID || "",
  );
  const guild = await oauthGuild.fetch();
  const emojis = await guild.emojis.fetch();
  const data = emojis.toJSON();
  return Response.json(data);
}
