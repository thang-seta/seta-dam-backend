import express from 'express';
import cors from 'cors';
import { ApolloServer } from '@apollo/server';
import { expressMiddleware } from '@apollo/server/express4';
import { ApolloServerPluginLandingPageLocalDefault } from '@apollo/server/plugin/landingPage/default';
import { typeDefs } from './graphql/schema';
import { resolvers } from './graphql/resolvers';
import { createContext, Context } from './graphql/context';

export async function createApp() {
  const app = express();

  // Set up standard middlewares
  app.use(cors());
  app.use(express.json());

  // Health check endpoint
  app.get('/health', (req, res) => {
    res.status(200).json({ status: 'OK' });
  });

  const isProd = process.env.NODE_ENV === 'production';
  const enableSandbox = process.env.ENABLE_SANDBOX === 'true' || !isProd;

  // Set up Apollo Server
  const server = new ApolloServer<Context>({
    typeDefs,
    resolvers,
    introspection: enableSandbox,
    plugins: enableSandbox
      ? [ApolloServerPluginLandingPageLocalDefault()]
      : [],
  });

  await server.start();

  app.use(
    '/graphql',
    expressMiddleware(server, {
      context: createContext,
    })
  );

  return app;
}
