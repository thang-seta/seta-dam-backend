import { IncomingMessage } from 'http';
import { UserContext } from '../clients/go-service.client';

const mockUsers: Record<string, UserContext> = {
  admin_user: {
    userId: '00000000-0000-0000-0000-000000000001',
    role: 'trainer_admin',
    username: 'admin_user',
  },
  editor_user: {
    userId: '00000000-0000-0000-0000-000000000002',
    role: 'editor',
    username: 'editor_user',
  },
  viewer_user: {
    userId: '00000000-0000-0000-0000-000000000003',
    role: 'viewer',
    username: 'viewer_user',
  },
};

export interface Context {
  user: UserContext | null;
}

export async function createContext({ req }: { req: IncomingMessage }): Promise<Context> {
  const userId = req.headers['x-user-id'] as string;
  const userRole = req.headers['x-user-role'] as string;
  const userUsername = req.headers['x-user-username'] as string;

  if (userId && userRole) {
    return {
      user: {
        userId,
        role: userRole,
        username: userUsername || 'unknown',
      },
    };
  }

  // Fallback: Check authorization header for mock user aliases
  const authHeader = req.headers['authorization'] || '';
  const username = authHeader.replace(/^Bearer\s+/i, '').trim();

  if (username && mockUsers[username]) {
    return {
      user: mockUsers[username],
    };
  }

  return {
    user: null,
  };
}
