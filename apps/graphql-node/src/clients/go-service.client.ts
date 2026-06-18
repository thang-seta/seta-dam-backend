import dotenv from 'dotenv';

dotenv.config();

const GO_SERVICE_URL = process.env.GO_SERVICE_URL || 'http://localhost:8080';

export interface UserContext {
  userId: string;
  role: string;
  username: string;
}

export class GoServiceClient {
  private static getHeaders(context: UserContext | null): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    if (context) {
      headers['X-User-ID'] = context.userId;
      headers['X-User-Role'] = context.role;
      headers['X-User-Username'] = context.username;
    }
    return headers;
  }

  private static async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      let errMsg = 'HTTP error';
      try {
        const errBody = await response.json() as any;
        errMsg = errBody.error || errBody.message || errMsg;
      } catch {
        errMsg = await response.text() || errMsg;
      }
      throw new Error(errMsg);
    }
    return response.json() as Promise<T>;
  }

  // Folder Operations
  static async getFolderTree(context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/folders/tree`, {
      method: 'GET',
      headers: this.getHeaders(context),
    });
    return this.handleResponse<any[]>(res);
  }

  static async getFolderById(id: string, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/folders/${id}`, {
      method: 'GET',
      headers: this.getHeaders(context),
    });
    return this.handleResponse<any>(res);
  }

  static async createFolder(name: string, description: string | null, parentId: string | null, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/folders`, {
      method: 'POST',
      headers: this.getHeaders(context),
      body: JSON.stringify({ name, description, parent_id: parentId }),
    });
    return this.handleResponse<any>(res);
  }

  static async updateFolder(id: string, name: string, description: string | null, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/folders/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(context),
      body: JSON.stringify({ name, description }),
    });
    return this.handleResponse<any>(res);
  }

  static async moveFolder(id: string, parentId: string | null, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/folders/${id}/move`, {
      method: 'PUT',
      headers: this.getHeaders(context),
      body: JSON.stringify({ parent_id: parentId }),
    });
    return this.handleResponse<any>(res);
  }

  static async deleteFolder(id: string, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/folders/${id}`, {
      method: 'DELETE',
      headers: this.getHeaders(context),
    });
    return this.handleResponse<any>(res);
  }

  // Metadata Operations
  static async getAllMetadata(context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/metadata`, {
      method: 'GET',
      headers: this.getHeaders(context),
    });
    return this.handleResponse<any[]>(res);
  }

  static async getMetadataList(folderId: string, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/folders/${folderId}/metadata`, {
      method: 'GET',
      headers: this.getHeaders(context),
    });
    return this.handleResponse<any[]>(res);
  }

  static async getMetadataById(id: string, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/metadata/${id}`, {
      method: 'GET',
      headers: this.getHeaders(context),
    });
    return this.handleResponse<any>(res);
  }

  static async createMetadata(meta: any, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/metadata`, {
      method: 'POST',
      headers: this.getHeaders(context),
      body: JSON.stringify(meta),
    });
    return this.handleResponse<any>(res);
  }

  static async updateMetadata(id: string, meta: any, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/metadata/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(context),
      body: JSON.stringify(meta),
    });
    return this.handleResponse<any>(res);
  }

  static async deleteMetadata(id: string, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/metadata/${id}`, {
      method: 'DELETE',
      headers: this.getHeaders(context),
    });
    return this.handleResponse<any>(res);
  }

  // Permission Operations
  static async grantPermission(perm: any, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/permissions`, {
      method: 'POST',
      headers: this.getHeaders(context),
      body: JSON.stringify(perm),
    });
    return this.handleResponse<any>(res);
  }

  static async revokePermission(perm: any, context: UserContext | null) {
    const res = await fetch(`${GO_SERVICE_URL}/api/permissions`, {
      method: 'DELETE',
      headers: this.getHeaders(context),
      body: JSON.stringify(perm),
    });
    return this.handleResponse<any>(res);
  }

  static async getEffectivePermissions(userId: string, objectType: string, objectId: string, context: UserContext | null) {
    const query = new URLSearchParams({
      user_id: userId,
      object_type: objectType,
      object_id: objectId,
    }).toString();
    const res = await fetch(`${GO_SERVICE_URL}/api/permissions/effective?${query}`, {
      method: 'GET',
      headers: this.getHeaders(context),
    });
    return this.handleResponse<any>(res);
  }
}
