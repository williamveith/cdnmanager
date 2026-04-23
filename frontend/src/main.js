import './styles/app.css';
import { initializeApp } from './controllers/appController';

window.addEventListener('DOMContentLoaded', async () => {
  await initializeApp();
});