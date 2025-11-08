export function LandingPage() {
  // 重定向到首页（模式选择页面）
  if (typeof window !== 'undefined') {
    window.location.href = '/';
    return null;
  }

  return null;
}
