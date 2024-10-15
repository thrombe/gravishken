"use client"
import { StrictMode, useEffect } from 'react'
import { createRoot } from 'react-dom/client';
import {
  createBrowserRouter,
  Outlet,
  RouterProvider,
  useNavigate,
} from "react-router-dom";
import './index.css';
import InstructionsPage from './pages/instructions';
import LoginPage from './pages/login-page';
import TestsPage from './pages/tests';
import EndPage from './pages/end';
import * as server from "@common/server.ts";
import * as types from "@common/types.ts";
import { TestProvider } from '@/components/TestContext';
import OfflineToast from './components/offline-toast';
import { useToast } from './hooks/use-toast';

function WebSocketHandler() {
  const navigate = useNavigate();

 const {toast} = useToast();

  useEffect(() => {
    let disable: (() => PromiseLike<void>)[] = [];

    server.server.add_callback(types.Varient.ExeNotFound, async (res) => {
      console.error('Error from server:', res.ErrMsg);
      toast({
        title: "Error",
        description: res.ErrMsg,
        variant:"destructive"
      })
    }).then(d => {
      disable.push(d);
    });

    server.server.add_callback(types.Varient.LoadRoute, async (res) => {
      console.log(res);
      navigate(res.Route)
    }).then(d => {
      disable.push(d);
    });

    server.server.add_callback(types.Varient.Err, async (res) => {
      console.error('Error from server:', res.Message);
      toast({
        title: "Error",
        description: res.Message,
        variant:"destructive"
      })
    }).then(d => {
      disable.push(d);
    });

    server.server.add_callback(types.Varient.WarnUser, async (res) => {
      console.error('Warning to the user:', res.Message);
      toast({
        title: "Warning",
        description: res.Message,
        variant:"destructive"
      })  
    }).then(d => {
      disable.push(d);
    });

    return () => {
      for (let fn of disable) {
        fn();
      }
    };
  }, [navigate]);

  return (
    <div>
      <Outlet />
    </div>
  );
}

const router = createBrowserRouter([
  {
    path: "/",
    element: <WebSocketHandler />,
    children: [
      {
        path: "/",
        element: <LoginPage />,
      },
      {
        path: "/instructions",
        element: <InstructionsPage />,
      },
      {
        path: "/tests",
        element: <TestProvider><TestsPage /></TestProvider>
      },
      {
        path: "/tests/:testId",
        element: <TestProvider><TestsPage /></TestProvider>
      },
      {
        path: "/end",
        element: <EndPage />
      }
    ],
  },
]);


server.init().then(async () => {
  createRoot(document.getElementById('root')!).render(
    <StrictMode>
      <RouterProvider router={router} />
      <OfflineToast/>
    </StrictMode>,
  );
});
