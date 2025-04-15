import React, { createContext, useContext, useEffect, useState } from 'react';
import * as authService from '../services/auth';

interface User {
  id: number;
  name: string;
  email: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (credentials: { email: string, password: string }) => Promise<boolean>;
  register: (name: string, email: string, password: string) => Promise<any>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const loadUser = async () => {
      const userData = await authService.getCurrentUser();
      if (userData) {
        setUser(userData);
      }
      setLoading(false);
    };
    
    loadUser();
  }, []);
  
  const handleLogin = async (credentials: { email: string, password: string }) => {
    const success = await authService.login(credentials);
    if (success) {
      const userData = await authService.getCurrentUser();
      setUser(userData);
    }
    return success;
  };
  
  const handleRegister = async (name: string, email: string, password: string) => {
    return await authService.register({ name, email, password });
  };
  
  const handleLogout = () => {
    authService.logout();
    setUser(null);
  };
  
  return (
    <AuthContext.Provider value={{
      user,
      loading,
      login: handleLogin,
      register: handleRegister,
      logout: handleLogout
    }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
