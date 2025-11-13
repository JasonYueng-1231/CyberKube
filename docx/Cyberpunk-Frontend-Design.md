# K8S-Manager 赛博朋克前端设计方案

**版本**: v1.0  
**创建日期**: 2025-11-11  
**设计风格**: 赛博朋克/未来科技风

---

## 目录

1. [设计参考](#1-设计参考)
2. [配色方案](#2-配色方案)
3. [核心组件设计](#3-核心组件设计)
4. [页面布局设计](#4-页面布局设计)
5. [动画效果](#5-动画效果)
6. [技术实现](#6-技术实现)
7. [组件库选择](#7-组件库选择)

---

## 1. 设计参考

### 1.1 参考界面分析（Kite界面）

基于提供的Kite界面截图，提取以下特点：

**布局结构**：
```
┌─────────────────────────────────────────────────────────┐
│  Sidebar                Top Bar                          │
│  ┌────────┐  ┌──────────────────────────────────────┐  │
│  │WORKLOADS│ │ 🌐集群选择  📦命名空间  🔍  🔔  👤     │  │
│  │        │  └──────────────────────────────────────┘  │
│  │• 概览  │                                             │
│  │• 部署  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐          │
│  │• Pod   │  │节点 │ │容器组│ │空间 │ │服务 │          │
│  │• 任务  │  │  8  │ │  1  │ │  5  │ │  6  │          │
│  │        │  └─────┘ └─────┘ └─────┘ └─────┘          │
│  │网络    │                                             │
│  │• 服务  │  ┌─────────────────────────────────────┐  │
│  │• Ingress│ │ CPU使用率                            │  │
│  │        │  │ ▓▓▓▓░░░░░░ 30% (Requests)          │  │
│  │配置    │  │ ▓▓▓▓▓▓▓▓▓▓ 200% (Limits) ⚠️       │  │
│  │• ConfigMap│                                      │  │
│  └────────┘  └─────────────────────────────────────┘  │
│              ┌─────────────────────────────────────┐  │
│              │ 最近事件                             │  │
│              │ No recent events                     │  │
│              └─────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

**需要增强的赛博元素**：
- ✨ 霓虹发光效果（边框、图标、文字）
- 🌊 数据流动画（统计卡片）
- 📡 扫描线效果（事件列表）
- 🔮 玻璃态背景（卡片、弹窗）
- ⚡ 悬停交互（表格行、按钮）
- 🎯 几何科技边框（面板边缘）
- 💫 粒子效果背景（可选）

---

## 2. 配色方案

### 2.1 核心色彩系统

```css
/* ==================== 背景色系 ==================== */
:root {
  /* 主背景 - 深空蓝黑 */
  --bg-primary: #0a0e27;
  --bg-secondary: #131729;
  --bg-tertiary: #1a1f3a;
  
  /* 卡片背景 - 半透明玻璃态 */
  --bg-card: rgba(26, 31, 58, 0.8);
  --bg-card-hover: rgba(36, 41, 70, 0.9);
  
  /* 侧边栏 */
  --bg-sidebar: rgba(19, 23, 41, 0.95);
  
  /* 弹窗/模态框 */
  --bg-modal: rgba(10, 14, 39, 0.98);
  --bg-modal-overlay: rgba(0, 0, 0, 0.7);
}

/* ==================== 霓虹色系 ==================== */
:root {
  /* 主色 - 青色霓虹 */
  --neon-cyan: #00f6ff;
  --neon-cyan-dark: #00b8c4;
  --neon-cyan-light: #5dffff;
  
  /* 辅色 - 紫色霓虹 */
  --neon-purple: #b537f2;
  --neon-purple-dark: #8b2ac0;
  --neon-purple-light: #d16eff;
  
  /* 品红霓虹 */
  --neon-magenta: #ff00ff;
  --neon-magenta-dark: #cc00cc;
  
  /* 绿色霓虹 - 成功状态 */
  --neon-green: #39ff14;
  --neon-green-glow: #2de600;
  
  /* 橙色霓虹 - 警告 */
  --neon-orange: #ff6600;
  --neon-orange-glow: #ff8833;
  
  /* 红色霓虹 - 错误/危险 */
  --neon-red: #ff006e;
  --neon-red-glow: #ff3388;
  
  /* 黄色霓虹 - 待处理 */
  --neon-yellow: #ffed4e;
}

/* ==================== 文字色系 ==================== */
:root {
  --text-primary: #e0e6ed;      /* 主要文字 */
  --text-secondary: #8b92a9;    /* 次要文字 */
  --text-muted: #4a5568;        /* 弱化文字 */
  --text-disabled: #2d3548;     /* 禁用文字 */
  
  /* 霓虹文字 */
  --text-neon-cyan: #00f6ff;
  --text-neon-green: #39ff14;
}

/* ==================== 边框/分割线 ==================== */
:root {
  --border-normal: #2d3548;
  --border-active: var(--neon-cyan);
  --border-glow: rgba(0, 246, 255, 0.3);
  --border-purple-glow: rgba(181, 55, 242, 0.3);
}

/* ==================== 阴影/发光效果 ==================== */
:root {
  /* 青色发光 */
  --shadow-cyan-sm: 0 0 10px rgba(0, 246, 255, 0.3);
  --shadow-cyan-md: 0 0 20px rgba(0, 246, 255, 0.5);
  --shadow-cyan-lg: 0 0 40px rgba(0, 246, 255, 0.7);
  
  /* 紫色发光 */
  --shadow-purple-sm: 0 0 10px rgba(181, 55, 242, 0.3);
  --shadow-purple-md: 0 0 20px rgba(181, 55, 242, 0.5);
  
  /* 绿色发光 */
  --shadow-green-sm: 0 0 10px rgba(57, 255, 20, 0.3);
  --shadow-green-md: 0 0 20px rgba(57, 255, 20, 0.5);
  
  /* 卡片阴影 */
  --shadow-card: 
    0 8px 32px rgba(0, 0, 0, 0.4),
    inset 0 0 20px rgba(0, 246, 255, 0.05);
}

/* ==================== 渐变色 ==================== */
:root {
  /* 背景渐变 */
  --gradient-bg: linear-gradient(135deg, #0a0e27 0%, #1a1f3a 100%);
  
  /* 霓虹渐变 */
  --gradient-cyan-purple: linear-gradient(135deg, #00f6ff 0%, #b537f2 100%);
  --gradient-cyan-blue: linear-gradient(135deg, #00f6ff 0%, #0080ff 100%);
  --gradient-purple-magenta: linear-gradient(135deg, #b537f2 0%, #ff00ff 100%);
  
  /* 状态渐变 */
  --gradient-success: linear-gradient(135deg, #39ff14 0%, #00ff88 100%);
  --gradient-warning: linear-gradient(135deg, #ffed4e 0%, #ff8800 100%);
  --gradient-danger: linear-gradient(135deg, #ff006e 0%, #ff4d00 100%);
  
  /* 卡片渐变 */
  --gradient-card: linear-gradient(145deg, #1a1f3a 0%, #131729 100%);
  
  /* 进度条渐变 */
  --gradient-progress-cyan: linear-gradient(90deg, 
    #00f6ff 0%, 
    rgba(0, 246, 255, 0.3) 100%);
  --gradient-progress-green: linear-gradient(90deg, 
    #39ff14 0%, 
    rgba(57, 255, 20, 0.3) 100%);
  --gradient-progress-red: linear-gradient(90deg, 
    #ff006e 0%, 
    rgba(255, 0, 110, 0.3) 100%);
}
```

### 2.2 状态色映射

```javascript
// 状态到颜色的映射
export const statusColors = {
  // Pod/Deployment 状态
  Running: { color: '#39ff14', glow: 'var(--shadow-green-md)' },
  Pending: { color: '#ffed4e', glow: '0 0 20px rgba(255, 237, 78, 0.5)' },
  Failed: { color: '#ff006e', glow: '0 0 20px rgba(255, 0, 110, 0.5)' },
  Succeeded: { color: '#00ff88', glow: '0 0 20px rgba(0, 255, 136, 0.5)' },
  Unknown: { color: '#8b92a9', glow: 'none' },
  
  // Node 状态
  Ready: { color: '#39ff14', glow: 'var(--shadow-green-md)' },
  NotReady: { color: '#ff006e', glow: '0 0 20px rgba(255, 0, 110, 0.5)' },
  
  // 健康状态
  Healthy: { color: '#00ff88', glow: 'var(--shadow-green-sm)' },
  Unhealthy: { color: '#ff6600', glow: '0 0 20px rgba(255, 102, 0, 0.5)' },
  Degraded: { color: '#ffed4e', glow: '0 0 20px rgba(255, 237, 78, 0.5)' },
};
```

---

## 3. 核心组件设计

### 3.1 统计卡片组件（Stats Card）

**设计要点**：
- 玻璃态背景
- 霓虹边框
- 图标发光效果
- 数据流动画
- 悬停立体效果

**组件代码**：

```tsx
// StatsCard.tsx
import React from 'react';
import './StatsCard.css';

interface StatsCardProps {
  icon: React.ReactNode;
  value: number | string;
  label: string;
  status?: string;
  color?: 'cyan' | 'purple' | 'green' | 'orange';
  trend?: 'up' | 'down' | 'stable';
}

export const StatsCard: React.FC<StatsCardProps> = ({
  icon,
  value,
  label,
  status,
  color = 'cyan',
  trend
}) => {
  return (
    <div className={`stats-card stats-card--${color}`}>
      {/* 数据流动画背景 */}
      <div className="stats-card__flow"></div>
      
      {/* 扫描线 */}
      <div className="stats-card__scan-line"></div>
      
      {/* 内容 */}
      <div className="stats-card__content">
        {/* 图标 - 带发光效果 */}
        <div className="stats-card__icon">
          {icon}
          <div className="icon-glow"></div>
        </div>
        
        {/* 数值 - 数字字体 + 渐变 */}
        <div className="stats-card__value">{value}</div>
        
        {/* 标签 */}
        <div className="stats-card__label">{label}</div>
        
        {/* 状态 */}
        {status && (
          <div className="stats-card__status">
            <span className="status-dot"></span>
            <span className="status-text">{status}</span>
          </div>
        )}
        
        {/* 趋势指示器（可选） */}
        {trend && (
          <div className={`stats-card__trend trend--${trend}`}>
            {trend === 'up' && '↑'}
            {trend === 'down' && '↓'}
            {trend === 'stable' && '→'}
          </div>
        )}
      </div>
      
      {/* 四角装饰 */}
      <div className="corner-decoration corner-tl"></div>
      <div className="corner-decoration corner-tr"></div>
      <div className="corner-decoration corner-bl"></div>
      <div className="corner-decoration corner-br"></div>
    </div>
  );
};
```

**样式代码**：

```css
/* StatsCard.css */
.stats-card {
  position: relative;
  padding: 24px;
  border-radius: 12px;
  background: var(--bg-card);
  backdrop-filter: blur(10px);
  border: 1px solid var(--border-glow);
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  min-height: 160px;
}

/* 悬停效果 - 立体感 */
.stats-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 
    var(--shadow-cyan-md),
    0 20px 40px rgba(0, 0, 0, 0.5);
  border-color: var(--neon-cyan);
}

/* 不同颜色主题 */
.stats-card--cyan {
  border-color: rgba(0, 246, 255, 0.3);
}
.stats-card--cyan:hover {
  box-shadow: 
    var(--shadow-cyan-md),
    0 20px 40px rgba(0, 0, 0, 0.5);
}

.stats-card--purple {
  border-color: rgba(181, 55, 242, 0.3);
}
.stats-card--purple:hover {
  box-shadow: 
    var(--shadow-purple-md),
    0 20px 40px rgba(0, 0, 0, 0.5);
}

.stats-card--green {
  border-color: rgba(57, 255, 20, 0.3);
}
.stats-card--green:hover {
  box-shadow: 
    var(--shadow-green-md),
    0 20px 40px rgba(0, 0, 0, 0.5);
}

/* ======== 数据流动画 ======== */
.stats-card__flow {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, 
    transparent 0%,
    rgba(0, 246, 255, 0.1) 50%,
    transparent 100%);
  animation: data-flow 3s ease-in-out infinite;
  pointer-events: none;
}

@keyframes data-flow {
  0% {
    left: -100%;
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    left: 100%;
    opacity: 0;
  }
}

/* ======== 扫描线 ======== */
.stats-card__scan-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, 
    transparent,
    var(--neon-cyan),
    transparent);
  opacity: 0.6;
  animation: scan-line 4s linear infinite;
  pointer-events: none;
}

@keyframes scan-line {
  0% {
    transform: translateY(0);
  }
  100% {
    transform: translateY(160px);
  }
}

/* ======== 内容区域 ======== */
.stats-card__content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* ======== 图标 ======== */
.stats-card__icon {
  position: relative;
  width: 48px;
  height: 48px;
  color: var(--neon-cyan);
  margin-bottom: 12px;
}

.stats-card--purple .stats-card__icon {
  color: var(--neon-purple);
}

.stats-card--green .stats-card__icon {
  color: var(--neon-green);
}

.stats-card__icon svg {
  width: 100%;
  height: 100%;
  filter: drop-shadow(0 0 8px currentColor);
  animation: icon-pulse 2s ease-in-out infinite;
}

@keyframes icon-pulse {
  0%, 100% {
    filter: drop-shadow(0 0 8px currentColor);
    opacity: 1;
  }
  50% {
    filter: drop-shadow(0 0 16px currentColor);
    opacity: 0.8;
  }
}

/* 图标发光效果 */
.icon-glow {
  position: absolute;
  inset: -10px;
  background: radial-gradient(circle, currentColor 0%, transparent 70%);
  opacity: 0.2;
  filter: blur(15px);
  animation: glow-pulse 2s ease-in-out infinite;
  pointer-events: none;
}

@keyframes glow-pulse {
  0%, 100% {
    opacity: 0.2;
  }
  50% {
    opacity: 0.4;
  }
}

/* ======== 数值 ======== */
.stats-card__value {
  font-size: 42px;
  font-weight: 700;
  font-family: 'Orbitron', 'Rajdhani', monospace;
  line-height: 1;
  background: var(--gradient-cyan-purple);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -1px;
}

/* ======== 标签 ======== */
.stats-card__label {
  font-size: 14px;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 600;
}

/* ======== 状态 ======== */
.stats-card__status {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-primary);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--neon-green);
  box-shadow: 0 0 10px var(--neon-green);
  animation: status-pulse 2s ease-in-out infinite;
}

@keyframes status-pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(1.2);
  }
}

/* ======== 四角装饰 ======== */
.corner-decoration {
  position: absolute;
  width: 12px;
  height: 12px;
  border: 2px solid var(--neon-cyan);
  opacity: 0.5;
  transition: opacity 0.3s ease;
}

.stats-card:hover .corner-decoration {
  opacity: 1;
}

.corner-tl {
  top: 8px;
  left: 8px;
  border-right: none;
  border-bottom: none;
}

.corner-tr {
  top: 8px;
  right: 8px;
  border-left: none;
  border-bottom: none;
}

.corner-bl {
  bottom: 8px;
  left: 8px;
  border-right: none;
  border-top: none;
}

.corner-br {
  bottom: 8px;
  right: 8px;
  border-left: none;
  border-top: none;
}
```

### 3.2 资源使用进度条（Progress Bar）

**设计要点**：
- 霓虹发光进度条
- 动态流光效果
- 超限警告动画
- 百分比悬停提示

**组件代码**：

```tsx
// ResourceProgress.tsx
interface ResourceProgressProps {
  label: string;
  value: number;
  max: number;
  unit?: string;
  type: 'requests' | 'limits';
  warningThreshold?: number; // 默认80%
}

export const ResourceProgress: React.FC<ResourceProgressProps> = ({
  label,
  value,
  max,
  unit = '',
  type,
  warningThreshold = 80
}) => {
  const percentage = (value / max) * 100;
  const isWarning = percentage >= warningThreshold;
  const isDanger = percentage >= 100;

  return (
    <div className="resource-progress">
      <div className="resource-progress__header">
        <span className="resource-progress__label">{label}</span>
        <span className="resource-progress__value">
          {value} {unit} / {max} {unit}
        </span>
      </div>

      <div className="resource-progress__track">
        <div 
          className={`resource-progress__fill ${
            isDanger ? 'fill--danger' : isWarning ? 'fill--warning' : 'fill--normal'
          }`}
          style={{ width: `${Math.min(percentage, 100)}%` }}
        >
          {/* 流光效果 */}
          <div className="progress-shine"></div>
          
          {/* 危险脉冲 */}
          {isDanger && <div className="danger-pulse"></div>}
        </div>
        
        {/* 网格背景 */}
        <div className="progress-grid"></div>
      </div>

      <div className={`resource-progress__percent ${
        isDanger ? 'text--danger' : isWarning ? 'text--warning' : ''
      }`}>
        {percentage.toFixed(1)}% of capacity
        {isDanger && <span className="warning-icon">⚠️</span>}
      </div>
    </div>
  );
};
```

**样式代码**：

```css
/* ResourceProgress.css */
.resource-progress {
  margin-bottom: 24px;
}

.resource-progress__header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 13px;
}

.resource-progress__label {
  color: var(--text-primary);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.resource-progress__value {
  color: var(--text-secondary);
  font-family: 'Fira Code', monospace;
}

/* ======== 进度条轨道 ======== */
.resource-progress__track {
  position: relative;
  height: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid rgba(0, 246, 255, 0.2);
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.3);
}

/* ======== 进度填充 ======== */
.resource-progress__fill {
  position: relative;
  height: 100%;
  border-radius: 6px;
  transition: width 1s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

/* 正常状态 - 青色 */
.fill--normal {
  background: var(--gradient-progress-cyan);
  box-shadow: 
    0 0 15px rgba(0, 246, 255, 0.5),
    inset 0 0 15px rgba(0, 246, 255, 0.3);
}

/* 警告状态 - 黄色 */
.fill--warning {
  background: linear-gradient(90deg, 
    #ffed4e 0%, 
    rgba(255, 237, 78, 0.3) 100%);
  box-shadow: 
    0 0 15px rgba(255, 237, 78, 0.5),
    inset 0 0 15px rgba(255, 237, 78, 0.3);
}

/* 危险状态 - 红色 */
.fill--danger {
  background: var(--gradient-progress-red);
  box-shadow: 
    0 0 15px rgba(255, 0, 110, 0.6),
    inset 0 0 15px rgba(255, 0, 110, 0.3);
}

/* ======== 流光效果 ======== */
.progress-shine {
  position: absolute;
  top: 0;
  left: -100%;
  width: 30%;
  height: 100%;
  background: linear-gradient(90deg, 
    transparent 0%,
    rgba(255, 255, 255, 0.4) 50%,
    transparent 100%);
  animation: shine 2s ease-in-out infinite;
}

@keyframes shine {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}

/* ======== 危险脉冲 ======== */
.danger-pulse {
  position: absolute;
  top: 0;
  right: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, 
    transparent 0%,
    rgba(255, 0, 110, 0.4) 100%);
  animation: pulse-danger 1s ease-in-out infinite;
}

@keyframes pulse-danger {
  0%, 100% {
    opacity: 0.3;
  }
  50% {
    opacity: 0.7;
  }
}

/* ======== 网格背景 ======== */
.progress-grid {
  position: absolute;
  inset: 0;
  background-image: 
    repeating-linear-gradient(
      90deg,
      rgba(255, 255, 255, 0.03) 0px,
      transparent 1px,
      transparent 10px
    );
  pointer-events: none;
}

/* ======== 百分比显示 ======== */
.resource-progress__percent {
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 6px;
}

.text--warning {
  color: var(--neon-yellow);
}

.text--danger {
  color: var(--neon-red);
  font-weight: 600;
}

.warning-icon {
  animation: warning-blink 1s ease-in-out infinite;
}

@keyframes warning-blink {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.3;
  }
}
```

### 3.3 霓虹按钮（Neon Button）

**多种样式**：
- Primary: 实心霓虹
- Secondary: 边框霓虹
- Ghost: 透明悬停发光
- Danger: 红色警告

**组件代码**：

```tsx
// NeonButton.tsx
interface NeonButtonProps {
  children: React.ReactNode;
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
  size?: 'small' | 'medium' | 'large';
  icon?: React.ReactNode;
  onClick?: () => void;
  disabled?: boolean;
  fullWidth?: boolean;
}

export const NeonButton: React.FC<NeonButtonProps> = ({
  children,
  variant = 'primary',
  size = 'medium',
  icon,
  onClick,
  disabled = false,
  fullWidth = false
}) => {
  return (
    <button
      className={`neon-btn neon-btn--${variant} neon-btn--${size} ${
        fullWidth ? 'neon-btn--full' : ''
      }`}
      onClick={onClick}
      disabled={disabled}
    >
      {/* 按钮背景层 */}
      <span className="btn-bg"></span>
      
      {/* 发光层 */}
      <span className="btn-glow"></span>
      
      {/* 边框动画 */}
      <span className="btn-border btn-border-top"></span>
      <span className="btn-border btn-border-right"></span>
      <span className="btn-border btn-border-bottom"></span>
      <span className="btn-border btn-border-left"></span>
      
      {/* 内容 */}
      <span className="btn-content">
        {icon && <span className="btn-icon">{icon}</span>}
        <span className="btn-text">{children}</span>
      </span>
    </button>
  );
};
```

**样式代码**：

```css
/* NeonButton.css */
.neon-btn {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12px 32px;
  border: none;
  background: transparent;
  color: var(--text-primary);
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.neon-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  pointer-events: none;
}

/* 按钮尺寸 */
.neon-btn--small {
  padding: 8px 20px;
  font-size: 12px;
}

.neon-btn--large {
  padding: 16px 40px;
  font-size: 16px;
}

.neon-btn--full {
  width: 100%;
}

/* ======== Primary 按钮 ======== */
.neon-btn--primary {
  color: #0a0e27;
}

.neon-btn--primary .btn-bg {
  position: absolute;
  inset: 0;
  background: var(--gradient-cyan-blue);
  z-index: 0;
}

.neon-btn--primary:hover .btn-bg {
  filter: brightness(1.2);
}

.neon-btn--primary .btn-glow {
  position: absolute;
  inset: -4px;
  background: var(--neon-cyan);
  filter: blur(20px);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: -1;
}

.neon-btn--primary:hover .btn-glow {
  opacity: 0.6;
  animation: glow-pulse-btn 2s ease-in-out infinite;
}

@keyframes glow-pulse-btn {
  0%, 100% {
    opacity: 0.4;
  }
  50% {
    opacity: 0.8;
  }
}

/* ======== Secondary 按钮 ======== */
.neon-btn--secondary {
  color: var(--neon-cyan);
}

.neon-btn--secondary::before {
  content: '';
  position: absolute;
  inset: 2px;
  background: var(--bg-primary);
  z-index: 0;
}

.neon-btn--secondary::after {
  content: '';
  position: absolute;
  inset: 0;
  border: 2px solid var(--neon-cyan);
  border-radius: inherit;
  box-shadow: 
    0 0 10px rgba(0, 246, 255, 0.3),
    inset 0 0 10px rgba(0, 246, 255, 0.1);
  transition: all 0.3s ease;
}

.neon-btn--secondary:hover::after {
  box-shadow: 
    0 0 20px rgba(0, 246, 255, 0.6),
    inset 0 0 20px rgba(0, 246, 255, 0.2);
}

.neon-btn--secondary:hover {
  background: rgba(0, 246, 255, 0.1);
}

/* ======== Ghost 按钮 ======== */
.neon-btn--ghost {
  color: var(--text-secondary);
}

.neon-btn--ghost:hover {
  color: var(--neon-cyan);
  background: rgba(0, 246, 255, 0.05);
}

/* ======== Danger 按钮 ======== */
.neon-btn--danger {
  color: #0a0e27;
}

.neon-btn--danger .btn-bg {
  position: absolute;
  inset: 0;
  background: var(--gradient-danger);
  z-index: 0;
}

.neon-btn--danger .btn-glow {
  position: absolute;
  inset: -4px;
  background: var(--neon-red);
  filter: blur(20px);
  opacity: 0;
  z-index: -1;
}

.neon-btn--danger:hover .btn-glow {
  opacity: 0.5;
}

/* ======== 边框动画 ======== */
.btn-border {
  position: absolute;
  background: var(--neon-cyan);
  opacity: 0;
  transition: all 0.3s ease;
}

.btn-border-top,
.btn-border-bottom {
  height: 2px;
  width: 0;
}

.btn-border-left,
.btn-border-right {
  width: 2px;
  height: 0;
}

.btn-border-top {
  top: 0;
  left: 0;
}

.btn-border-right {
  top: 0;
  right: 0;
}

.btn-border-bottom {
  bottom: 0;
  right: 0;
}

.btn-border-left {
  bottom: 0;
  left: 0;
}

/* Primary按钮悬停边框 */
.neon-btn--primary:hover .btn-border {
  opacity: 1;
}

.neon-btn--primary:hover .btn-border-top {
  width: 100%;
  transition-delay: 0s;
}

.neon-btn--primary:hover .btn-border-right {
  height: 100%;
  transition-delay: 0.15s;
}

.neon-btn--primary:hover .btn-border-bottom {
  width: 100%;
  transition-delay: 0.3s;
}

.neon-btn--primary:hover .btn-border-left {
  height: 100%;
  transition-delay: 0.45s;
}

/* ======== 按钮内容 ======== */
.btn-content {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-icon svg {
  width: 16px;
  height: 16px;
}

/* 点击涟漪效果 */
.neon-btn::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  transform: translate(-50%, -50%);
  transition: width 0.6s ease, height 0.6s ease;
  pointer-events: none;
}

.neon-btn:active::before {
  width: 300px;
  height: 300px;
  transition: width 0s, height 0s;
}
```

### 3.4 数据表格（Cyber Table）

**特点**：
- 悬停行发光
- 左侧光带
- 扫描线动画
- 状态徽章
- 操作按钮悬浮

**组件代码**：

```tsx
// CyberTable.tsx
interface Column {
  key: string;
  title: string;
  width?: string;
  render?: (value: any, record: any) => React.ReactNode;
}

interface CyberTableProps {
  columns: Column[];
  dataSource: any[];
  onRowClick?: (record: any) => void;
}

export const CyberTable: React.FC<CyberTableProps> = ({
  columns,
  dataSource,
  onRowClick
}) => {
  return (
    <div className="cyber-table-wrapper">
      {/* 扫描线背景 */}
      <div className="table-scan-effect"></div>
      
      <table className="cyber-table">
        <thead>
          <tr className="table-header-row">
            {columns.map(col => (
              <th 
                key={col.key}
                style={{ width: col.width }}
              >
                <div className="table-header-cell">
                  {col.title}
                  <div className="header-glow"></div>
                </div>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {dataSource.map((record, index) => (
            <tr 
              key={index}
              className="table-row"
              onClick={() => onRowClick?.(record)}
            >
              {/* 左侧光带 */}
              <div className="row-light-strip"></div>
              
              {columns.map(col => (
                <td key={col.key}>
                  <div className="table-cell">
                    {col.render 
                      ? col.render(record[col.key], record)
                      : record[col.key]
                    }
                  </div>
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

// 状态徽章组件
export const StatusBadge: React.FC<{
  status: string;
  color?: string;
}> = ({ status, color }) => {
  return (
    <div className={`status-badge status-badge--${color || 'default'}`}>
      <span className="status-dot"></span>
      <span className="status-text">{status}</span>
    </div>
  );
};
```

**样式代码**：

```css
/* CyberTable.css */
.cyber-table-wrapper {
  position: relative;
  background: var(--bg-card);
  border-radius: 12px;
  border: 1px solid var(--border-glow);
  overflow: hidden;
  backdrop-filter: blur(10px);
}

/* ======== 扫描线背景 ======== */
.table-scan-effect {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, 
    transparent,
    var(--neon-cyan),
    transparent);
  opacity: 0.6;
  animation: table-scan 5s linear infinite;
  pointer-events: none;
  z-index: 10;
}

@keyframes table-scan {
  0% {
    transform: translateY(0);
  }
  100% {
    transform: translateY(600px);
  }
}

/* ======== 表格 ======== */
.cyber-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
}

/* ======== 表头 ======== */
.table-header-row {
  background: rgba(0, 246, 255, 0.05);
  border-bottom: 2px solid var(--border-glow);
}

.table-header-row th {
  padding: 16px;
  text-align: left;
  font-weight: 600;
  color: var(--neon-cyan);
  text-transform: uppercase;
  font-size: 12px;
  letter-spacing: 1.5px;
  position: relative;
}

.table-header-cell {
  position: relative;
  display: inline-block;
}

.header-glow {
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 0;
  height: 2px;
  background: var(--neon-cyan);
  box-shadow: 0 0 10px var(--neon-cyan);
  transition: width 0.3s ease;
}

.table-header-row th:hover .header-glow {
  width: 100%;
}

/* ======== 表格行 ======== */
.table-row {
  position: relative;
  background: transparent;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.table-row::before {
  content: '';
  position: absolute;
  inset: 0;
  background: rgba(0, 246, 255, 0.02);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.table-row:hover::before {
  opacity: 1;
}

.table-row:hover {
  background: rgba(0, 246, 255, 0.05);
  box-shadow: 
    inset 0 0 20px rgba(0, 246, 255, 0.1),
    0 0 10px rgba(0, 246, 255, 0.2);
  transform: translateX(8px);
}

/* 左侧光带 */
.row-light-strip {
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%;
  background: var(--gradient-cyan-purple);
  opacity: 0;
  box-shadow: 0 0 15px var(--neon-cyan);
  transition: opacity 0.3s ease;
}

.table-row:hover .row-light-strip {
  opacity: 1;
  animation: light-strip-pulse 2s ease-in-out infinite;
}

@keyframes light-strip-pulse {
  0%, 100% {
    box-shadow: 0 0 10px var(--neon-cyan);
  }
  50% {
    box-shadow: 0 0 20px var(--neon-cyan);
  }
}

/* ======== 表格单元格 ======== */
.table-row td {
  padding: 16px;
  color: var(--text-primary);
  font-size: 14px;
}

.table-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* ======== 状态徽章 ======== */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  border: 1px solid;
}

.status-badge--default {
  background: rgba(139, 146, 169, 0.1);
  border-color: rgba(139, 146, 169, 0.3);
  color: var(--text-secondary);
}

.status-badge--success {
  background: rgba(57, 255, 20, 0.1);
  border-color: rgba(57, 255, 20, 0.3);
  color: var(--neon-green);
}

.status-badge--warning {
  background: rgba(255, 237, 78, 0.1);
  border-color: rgba(255, 237, 78, 0.3);
  color: var(--neon-yellow);
}

.status-badge--danger {
  background: rgba(255, 0, 110, 0.1);
  border-color: rgba(255, 0, 110, 0.3);
  color: var(--neon-red);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 10px currentColor;
  animation: status-pulse 2s ease-in-out infinite;
}

@keyframes status-pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
```

---

## 4. 页面布局实现

### 4.1 Dashboard 页面完整代码

```tsx
// pages/Dashboard.tsx
import React from 'react';
import { StatsCard } from '../components/StatsCard';
import { ResourceProgress } from '../components/ResourceProgress';
import {
  ServerIcon,
  CubeIcon,
  FolderIcon,
  GlobeIcon,
  CpuChipIcon,
  CircleStackIcon
} from '@heroicons/react/24/outline';
import './Dashboard.css';

export const Dashboard: React.FC = () => {
  return (
    <div className="dashboard">
      {/* 背景效果 */}
      <div className="dashboard-bg">
        <div className="grid-pattern"></div>
        <div className="gradient-orb orb-1"></div>
        <div className="gradient-orb orb-2"></div>
      </div>

      <div className="dashboard-content">
        {/* 顶部统计卡片 */}
        <section className="stats-section">
          <div className="stats-grid">
            <StatsCard
              icon={<ServerIcon />}
              value="8"
              label="节点"
              status="All ready"
              color="cyan"
            />
            <StatsCard
              icon={<CubeIcon />}
              value="1"
              label="容器组"
              status="All ready"
              color="purple"
            />
            <StatsCard
              icon={<FolderIcon />}
              value="5"
              label="命名空间"
              status="All ready"
              color="green"
            />
            <StatsCard
              icon={<GlobeIcon />}
              value="6"
              label="服务"
              status="All ready"
              color="orange"
            />
          </div>
        </section>

        {/* 资源使用情况 */}
        <section className="resource-section">
          <div className="section-grid">
            {/* CPU使用率 */}
            <div className="resource-card glass-card">
              <div className="card-header">
                <h3 className="card-title">
                  <CpuChipIcon className="title-icon" />
                  CPU使用率
                </h3>
              </div>
              <div className="card-body">
                <div className="resource-summary">
                  Requests: 0.6 / Limits: 4.0 / Total: 2.00 cores
                </div>
                <div className="progress-group">
                  <ResourceProgress
                    label="Requests"
                    value={0.6}
                    max={2.0}
                    unit="cores"
                    type="requests"
                  />
                  <ResourceProgress
                    label="Limits"
                    value={4.0}
                    max={2.0}
                    unit="cores"
                    type="limits"
                    warningThreshold={80}
                  />
                </div>
                <div className="resource-footer">
                  <span className="footer-label">Available:</span>
                  <span className="footer-value neon-text">1.4 cores</span>
                </div>
              </div>
            </div>

            {/* 内存使用率 */}
            <div className="resource-card glass-card">
              <div className="card-header">
                <h3 className="card-title">
                  <CircleStackIcon className="title-icon" />
                  内存使用率
                </h3>
              </div>
              <div className="card-body">
                <div className="resource-summary">
                  Requests: 0.3 / Limits: 1.5 / Total: 1.88 GiB
                </div>
                <div className="progress-group">
                  <ResourceProgress
                    label="Requests"
                    value={0.3}
                    max={1.88}
                    unit="GiB"
                    type="requests"
                  />
                  <ResourceProgress
                    label="Limits"
                    value={1.5}
                    max={1.88}
                    unit="GiB"
                    type="limits"
                    warningThreshold={80}
                  />
                </div>
                <div className="resource-footer">
                  <span className="footer-label">Available:</span>
                  <span className="footer-value neon-text">1.6 GiB</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* 最近事件 */}
        <section className="events-section">
          <div className="events-card glass-card">
            <div className="card-header">
              <h3 className="card-title">
                <span className="title-icon">📡</span>
                最近事件
              </h3>
              <div className="scan-line"></div>
            </div>
            <div className="events-empty">
              <div className="hologram-icon">
                <div className="hologram-rings">
                  <div className="ring ring-1"></div>
                  <div className="ring ring-2"></div>
                  <div className="ring ring-3"></div>
                </div>
                <div className="icon-center">✓</div>
              </div>
              <p className="empty-text">No recent events</p>
              <p className="empty-subtext">系统运行正常</p>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
};
```

**样式代码**：

```css
/* Dashboard.css */
.dashboard {
  position: relative;
  min-height: 100vh;
  padding: 24px;
  background: var(--bg-primary);
}

/* ======== 背景效果 ======== */
.dashboard-bg {
  position: fixed;
  inset: 0;
  z-index: 0;
  overflow: hidden;
  pointer-events: none;
}

/* 网格背景 */
.grid-pattern {
  position: absolute;
  inset: 0;
  background-image: 
    linear-gradient(rgba(0, 246, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 246, 255, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: grid-flow 20s linear infinite;
}

@keyframes grid-flow {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 50px 50px;
  }
}

/* 渐变光球 */
.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.3;
  animation: orb-float 10s ease-in-out infinite;
}

.orb-1 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, var(--neon-cyan) 0%, transparent 70%);
  top: -200px;
  left: -200px;
}

.orb-2 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, var(--neon-purple) 0%, transparent 70%);
  bottom: -250px;
  right: -250px;
  animation-delay: -5s;
}

@keyframes orb-float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(20px, 20px) scale(1.1);
  }
}

/* ======== 内容区域 ======== */
.dashboard-content {
  position: relative;
  z-index: 1;
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ======== 统计卡片区 ======== */
.stats-section {
  animation: fade-in-up 0.6s ease-out;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ======== 资源使用区 ======== */
.resource-section {
  animation: fade-in-up 0.6s ease-out 0.2s;
  animation-fill-mode: both;
}

.section-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
}

.resource-card {
  padding: 24px;
  border-radius: 12px;
  background: var(--bg-card);
  backdrop-filter: blur(10px);
  border: 1px solid var(--border-glow);
  transition: all 0.3s ease;
}

.resource-card:hover {
  border-color: var(--neon-cyan);
  box-shadow: 
    var(--shadow-cyan-sm),
    0 8px 24px rgba(0, 0, 0, 0.3);
}

.card-header {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(0, 246, 255, 0.1);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.title-icon {
  width: 24px;
  height: 24px;
  color: var(--neon-cyan);
  filter: drop-shadow(0 0 8px var(--neon-cyan));
}

.resource-summary {
  margin-bottom: 20px;
  padding: 12px;
  background: rgba(0, 246, 255, 0.05);
  border-left: 3px solid var(--neon-cyan);
  border-radius: 4px;
  font-size: 13px;
  color: var(--text-secondary);
  font-family: 'Fira Code', monospace;
}

.progress-group {
  margin-bottom: 16px;
}

.resource-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 246, 255, 0.1);
  font-size: 14px;
}

.footer-label {
  color: var(--text-secondary);
}

.footer-value {
  font-weight: 600;
  font-family: 'Orbitron', monospace;
}

.neon-text {
  color: var(--neon-cyan);
  text-shadow: 0 0 10px var(--neon-cyan);
}

/* ======== 事件区域 ======== */
.events-section {
  animation: fade-in-up 0.6s ease-out 0.4s;
  animation-fill-mode: both;
}

.events-card {
  padding: 24px;
  border-radius: 12px;
  background: var(--bg-card);
  backdrop-filter: blur(10px);
  border: 1px solid var(--border-glow);
  min-height: 250px;
  position: relative;
  overflow: hidden;
}

/* 扫描线 */
.scan-line {
  position: absolute;
  top: 50px;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, 
    transparent,
    var(--neon-cyan),
    transparent);
  opacity: 0.6;
  animation: scan-horizontal 5s linear infinite;
}

@keyframes scan-horizontal {
  0% {
    transform: translateY(0);
  }
  100% {
    transform: translateY(200px);
  }
}

/* 空状态 */
.events-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  position: relative;
  z-index: 1;
}

/* 全息图标 */
.hologram-icon {
  position: relative;
  width: 80px;
  height: 80px;
  margin-bottom: 20px;
}

.hologram-rings {
  position: absolute;
  inset: 0;
}

.ring {
  position: absolute;
  inset: 0;
  border: 2px solid var(--neon-cyan);
  border-radius: 50%;
  animation: ring-pulse 3s ease-in-out infinite;
}

.ring-1 {
  animation-delay: 0s;
}

.ring-2 {
  animation-delay: 1s;
}

.ring-3 {
  animation-delay: 2s;
}

@keyframes ring-pulse {
  0% {
    transform: scale(0.8);
    opacity: 1;
  }
  100% {
    transform: scale(1.2);
    opacity: 0;
  }
}

.icon-center {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  color: var(--neon-green);
  text-shadow: 0 0 20px var(--neon-green);
}

.empty-text {
  font-size: 16px;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.empty-subtext {
  font-size: 14px;
  color: var(--text-muted);
}

/* ======== 响应式 ======== */
@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .section-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .dashboard {
    padding: 16px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .dashboard-content {
    gap: 16px;
  }
}
```

---

## 5. 动画效果库

### 5.1 全局动画

```css
/* animations.css - 全局动画库 */

/* ==================== 入场动画 ==================== */
@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fade-in-down {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slide-in-left {
  from {
    opacity: 0;
    transform: translateX(-50px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes slide-in-right {
  from {
    opacity: 0;
    transform: translateX(50px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes scale-in {
  from {
    opacity: 0;
    transform: scale(0.8);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* ==================== 发光/脉冲动画 ==================== */
@keyframes neon-pulse {
  0%, 100% {
    opacity: 1;
    box-shadow: 
      0 0 10px currentColor,
      0 0 20px currentColor;
  }
  50% {
    opacity: 0.7;
    box-shadow: 
      0 0 20px currentColor,
      0 0 40px currentColor;
  }
}

@keyframes glow-pulse {
  0%, 100% {
    filter: drop-shadow(0 0 8px currentColor);
  }
  50% {
    filter: drop-shadow(0 0 16px currentColor);
  }
}

@keyframes border-glow {
  0%, 100% {
    box-shadow: 
      0 0 5px var(--neon-cyan),
      0 0 10px var(--neon-cyan);
  }
  50% {
    box-shadow: 
      0 0 10px var(--neon-cyan),
      0 0 20px var(--neon-cyan),
      0 0 30px var(--neon-cyan);
  }
}

/* ==================== 数据流动画 ==================== */
@keyframes data-stream {
  0% {
    transform: translateX(-100%);
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    transform: translateX(100%);
    opacity: 0;
  }
}

@keyframes data-flow-vertical {
  0% {
    transform: translateY(-100%);
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    transform: translateY(100%);
    opacity: 0;
  }
}

/* ==================== 扫描线动画 ==================== */
@keyframes scan-line {
  0% {
    transform: translateY(0);
  }
  100% {
    transform: translateY(500px);
  }
}

@keyframes scan-horizontal {
  0% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(100%);
  }
}

/* ==================== 旋转动画 ==================== */
@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@keyframes rotate-reverse {
  from {
    transform: rotate(360deg);
  }
  to {
    transform: rotate(0deg);
  }
}

/* ==================== 波纹动画 ==================== */
@keyframes ripple {
  0% {
    transform: scale(0);
    opacity: 1;
  }
  100% {
    transform: scale(4);
    opacity: 0;
  }
}

/* ==================== 闪烁动画 ==================== */
@keyframes blink {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.3;
  }
}

@keyframes flicker {
  0%, 100% {
    opacity: 1;
  }
  10% {
    opacity: 0.8;
  }
  20% {
    opacity: 1;
  }
  30% {
    opacity: 0.9;
  }
  40% {
    opacity: 1;
  }
}

/* ==================== 浮动动画 ==================== */
@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

@keyframes float-horizontal {
  0%, 100% {
    transform: translateX(0);
  }
  50% {
    transform: translateX(10px);
  }
}

/* ==================== 数字滚动动画 ==================== */
@keyframes number-roll {
  0% {
    transform: translateY(100%);
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    transform: translateY(0);
  }
}

/* ==================== 工具类 ==================== */
.animate-fade-in {
  animation: fade-in 0.5s ease-out;
}

.animate-fade-in-up {
  animation: fade-in-up 0.6s ease-out;
}

.animate-slide-in-left {
  animation: slide-in-left 0.5s ease-out;
}

.animate-scale-in {
  animation: scale-in 0.4s ease-out;
}

.animate-neon-pulse {
  animation: neon-pulse 2s ease-in-out infinite;
}

.animate-glow-pulse {
  animation: glow-pulse 2s ease-in-out infinite;
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}

.animate-rotate {
  animation: rotate 20s linear infinite;
}

/* 延迟工具类 */
.delay-100 {
  animation-delay: 0.1s;
}

.delay-200 {
  animation-delay: 0.2s;
}

.delay-300 {
  animation-delay: 0.3s;
}

.delay-400 {
  animation-delay: 0.4s;
}

.delay-500 {
  animation-delay: 0.5s;
}
```

---

## 6. 技术实现

### 6.1 技术栈

**前端框架**：
- React 18+
- TypeScript 5+
- Vite (构建工具)

**UI库**：
- Ant Design 5+ (基础组件)
- 自定义赛博风格覆盖层

**动画库**：
- Framer Motion (高级动画)
- 原生CSS动画 (性能更好)

**图表库**：
- ECharts 5+ (数据可视化)
- 自定义赛博主题

**工具库**：
- Axios (HTTP请求)
- React Query (数据缓存)
- Zustand (状态管理)
- Day.js (时间处理)

### 6.2 项目结构

```
frontend/
├── src/
│   ├── assets/
│   │   ├── styles/
│   │   │   ├── global.css           # 全局样式
│   │   │   ├── variables.css        # CSS变量
│   │   │   ├── animations.css       # 动画库
│   │   │   └── cyberpunk.css        # 赛博风格
│   │   └── fonts/                   # 字体文件
│   │       ├── Orbitron/
│   │       ├── Rajdhani/
│   │       └── FiraCode/
│   │
│   ├── components/
│   │   ├── common/                  # 通用组件
│   │   │   ├── NeonButton/
│   │   │   ├── StatsCard/
│   │   │   ├── GlassCard/
│   │   │   ├── CyberTable/
│   │   │   ├── StatusBadge/
│   │   │   ├── ResourceProgress/
│   │   │   └── LoadingSpinner/
│   │   │
│   │   ├── layout/                  # 布局组件
│   │   │   ├── Header/
│   │   │   ├── Sidebar/
│   │   │   └── Footer/
│   │   │
│   │   └── business/                # 业务组件
│   │       ├── ClusterSelector/
│   │       ├── NamespaceSelector/
│   │       ├── PodList/
│   │       └── DeploymentCard/
│   │
│   ├── pages/
│   │   ├── Dashboard/               # 仪表盘
│   │   ├── Cluster/                 # 集群管理
│   │   ├── Workload/                # 工作负载
│   │   │   ├── Deployment/
│   │   │   ├── Pod/
│   │   │   └── StatefulSet/
│   │   ├── Service/                 # 服务
│   │   ├── Config/                  # 配置
│   │   └── Node/                    # 节点
│   │
│   ├── hooks/                       # 自定义Hooks
│   │   ├── useCluster.ts
│   │   ├── useDeployment.ts
│   │   └── useWebSocket.ts
│   │
│   ├── services/                    # API服务
│   │   ├── api.ts
│   │   ├── cluster.ts
│   │   └── deployment.ts
│   │
│   ├── store/                       # 状态管理
│   │   ├── useClusterStore.ts
│   │   └── useUserStore.ts
│   │
│   ├── utils/                       # 工具函数
│   │   ├── format.ts
│   │   └── constants.ts
│   │
│   ├── types/                       # TypeScript类型
│   │   ├── cluster.ts
│   │   └── deployment.ts
│   │
│   ├── App.tsx
│   └── main.tsx
│
├── public/
├── package.json
├── tsconfig.json
└── vite.config.ts
```

### 6.3 Ant Design自定义主题

```typescript
// theme.ts - Ant Design 赛博朋克主题配置
import type { ThemeConfig } from 'antd';

export const cyberTheme: ThemeConfig = {
  token: {
    // 主色
    colorPrimary: '#00f6ff',
    colorSuccess: '#39ff14',
    colorWarning: '#ffed4e',
    colorError: '#ff006e',
    colorInfo: '#0080ff',
    
    // 背景色
    colorBgBase: '#0a0e27',
    colorBgContainer: '#1a1f3a',
    colorBgElevated: '#242946',
    
    // 文字色
    colorText: '#e0e6ed',
    colorTextSecondary: '#8b92a9',
    colorTextTertiary: '#4a5568',
    colorTextQuaternary: '#2d3548',
    
    // 边框
    colorBorder: '#2d3548',
    colorBorderSecondary: 'rgba(0, 246, 255, 0.2)',
    
    // 字体
    fontFamily: "'Rajdhani', 'Inter', -apple-system, BlinkMacSystemFont, sans-serif",
    fontFamilyCode: "'Fira Code', 'JetBrains Mono', monospace",
    
    // 圆角
    borderRadius: 8,
    borderRadiusLG: 12,
    borderRadiusSM: 6,
    
    // 阴影
    boxShadow: '0 8px 32px rgba(0, 0, 0, 0.4)',
    boxShadowSecondary: '0 4px 16px rgba(0, 0, 0, 0.3)',
  },
  
  components: {
    Button: {
      primaryShadow: '0 0 20px rgba(0, 246, 255, 0.5)',
      dangerShadow: '0 0 20px rgba(255, 0, 110, 0.5)',
    },
    
    Table: {
      headerBg: 'rgba(0, 246, 255, 0.05)',
      headerColor: '#00f6ff',
      rowHoverBg: 'rgba(0, 246, 255, 0.05)',
      borderColor: 'rgba(0, 246, 255, 0.2)',
    },
    
    Card: {
      colorBgContainer: 'rgba(26, 31, 58, 0.8)',
      boxShadowTertiary: '0 8px 32px rgba(0, 0, 0, 0.4)',
    },
    
    Input: {
      colorBgContainer: 'rgba(19, 23, 41, 0.6)',
      colorBorder: 'rgba(0, 246, 255, 0.3)',
      activeBorderColor: '#00f6ff',
      hoverBorderColor: '#00f6ff',
      activeShadow: '0 0 10px rgba(0, 246, 255, 0.3)',
    },
    
    Select: {
      colorBgContainer: 'rgba(19, 23, 41, 0.6)',
      colorBorder: 'rgba(0, 246, 255, 0.3)',
      colorBgElevated: '#1a1f3a',
      optionSelectedBg: 'rgba(0, 246, 255, 0.1)',
    },
    
    Modal: {
      contentBg: 'rgba(26, 31, 58, 0.95)',
      headerBg: 'rgba(26, 31, 58, 0.95)',
      footerBg: 'rgba(26, 31, 58, 0.95)',
    },
  },
  
  algorithm: 'dark', // 使用暗色算法
};
```

**使用方式**：

```tsx
// App.tsx
import { ConfigProvider } from 'antd';
import { cyberTheme } from './theme';

function App() {
  return (
    <ConfigProvider theme={cyberTheme}>
      {/* 你的应用 */}
    </ConfigProvider>
  );
}
```

### 6.4 全局样式设置

```css
/* global.css */
@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;500;600;700;800;900&family=Rajdhani:wght@300;400;500;600;700&display=swap');
@import url('https://fonts.googleapis.com/css2?family=Fira+Code:wght@300;400;500;600;700&display=swap');

/* 重置和基础样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  width: 100%;
  height: 100%;
  overflow-x: hidden;
}

body {
  font-family: 'Rajdhani', 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 16px;
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, var(--neon-cyan), var(--neon-purple));
  border-radius: 10px;
  box-shadow: 0 0 10px var(--neon-cyan);
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, var(--neon-cyan-light), var(--neon-purple-light));
}

/* 选中文字样式 */
::selection {
  background: rgba(0, 246, 255, 0.3);
  color: var(--text-primary);
}

/* 禁用文字选择的区域 */
.no-select {
  user-select: none;
  -webkit-user-select: none;
}

/* 链接样式 */
a {
  color: var(--neon-cyan);
  text-decoration: none;
  transition: all 0.3s ease;
}

a:hover {
  color: var(--neon-cyan-light);
  text-shadow: 0 0 10px var(--neon-cyan);
}

/* 代码样式 */
code {
  font-family: 'Fira Code', 'JetBrains Mono', monospace;
  background: rgba(0, 246, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.9em;
  color: var(--neon-cyan);
}

/* 标题样式 */
h1, h2, h3, h4, h5, h6 {
  font-family: 'Orbitron', sans-serif;
  font-weight: 600;
  line-height: 1.2;
  color: var(--text-primary);
}

h1 {
  font-size: 2.5rem;
}

h2 {
  font-size: 2rem;
}

h3 {
  font-size: 1.5rem;
}

/* 按钮重置 */
button {
  font-family: inherit;
  cursor: pointer;
  border: none;
  background: none;
  outline: none;
}

button:focus-visible {
  outline: 2px solid var(--neon-cyan);
  outline-offset: 2px;
}

/* 输入框重置 */
input,
textarea,
select {
  font-family: inherit;
  color: inherit;
  background: transparent;
  border: none;
  outline: none;
}

input:focus,
textarea:focus,
select:focus {
  outline: none;
}

/* 列表样式 */
ul,
ol {
  list-style: none;
}

/* 图片样式 */
img {
  max-width: 100%;
  height: auto;
  display: block;
}

/* 加载动画 */
.loading {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 3px solid rgba(0, 246, 255, 0.3);
  border-top-color: var(--neon-cyan);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 工具类 */
.text-center {
  text-align: center;
}

.text-left {
  text-align: left;
}

.text-right {
  text-align: right;
}

.flex {
  display: flex;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.grid {
  display: grid;
}

/* 间距工具类 */
.mt-1 { margin-top: 0.5rem; }
.mt-2 { margin-top: 1rem; }
.mt-3 { margin-top: 1.5rem; }
.mt-4 { margin-top: 2rem; }

.mb-1 { margin-bottom: 0.5rem; }
.mb-2 { margin-bottom: 1rem; }
.mb-3 { margin-bottom: 1.5rem; }
.mb-4 { margin-bottom: 2rem; }

.ml-1 { margin-left: 0.5rem; }
.ml-2 { margin-left: 1rem; }
.ml-3 { margin-left: 1.5rem; }
.ml-4 { margin-left: 2rem; }

.mr-1 { margin-right: 0.5rem; }
.mr-2 { margin-right: 1rem; }
.mr-3 { margin-right: 1.5rem; }
.mr-4 { margin-right: 2rem; }

.p-1 { padding: 0.5rem; }
.p-2 { padding: 1rem; }
.p-3 { padding: 1.5rem; }
.p-4 { padding: 2rem; }
```

---

## 7. 组件库整合

### 7.1 推荐使用的开源组件

**基础UI**:
- Ant Design 5+ (主力UI库)
- Heroicons (图标库)

**图表可视化**:
- ECharts (推荐，功能强大)
- Recharts (React友好)

**动画**:
- Framer Motion (高级动画)
- React Spring (物理动画)

**终端/编辑器**:
- xterm.js (Web终端)
- Monaco Editor (代码编辑器)

**工具库**:
- clsx (className管理)
- dayjs (时间处理)
- lodash (工具函数)

### 7.2 完整的package.json

```json
{
  "name": "k8s-manager-frontend",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.20.0",
    
    "antd": "^5.12.0",
    "@ant-design/icons": "^5.2.6",
    
    "axios": "^1.6.2",
    "@tanstack/react-query": "^5.12.0",
    
    "zustand": "^4.4.7",
    
    "echarts": "^5.4.3",
    "echarts-for-react": "^3.0.2",
    
    "framer-motion": "^10.16.15",
    
    "xterm": "^5.3.0",
    "xterm-addon-fit": "^0.8.0",
    "xterm-addon-web-links": "^0.9.0",
    
    "@monaco-editor/react": "^4.6.0",
    
    "dayjs": "^1.11.10",
    "lodash-es": "^4.17.21",
    "clsx": "^2.0.0",
    
    "@heroicons/react": "^2.1.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.43",
    "@types/react-dom": "^18.2.17",
    "@types/lodash-es": "^4.17.12",
    
    "@typescript-eslint/eslint-plugin": "^6.14.0",
    "@typescript-eslint/parser": "^6.14.0",
    
    "@vitejs/plugin-react": "^4.2.1",
    "typescript": "^5.2.2",
    "vite": "^5.0.8",
    
    "eslint": "^8.55.0",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-react-refresh": "^0.4.5"
  }
}
```

---

## 8. 总结

### 8.1 核心特点

✨ **视觉效果**:
- 深色赛博朋克主题
- 霓虹发光效果
- 玻璃态卡片
- 流光数据动画
- 扫描线特效

⚡ **交互体验**:
- 流畅的过渡动画
- 悬停立体效果
- 点击反馈
- 加载状态

🎨 **设计系统**:
- 统一的配色方案
- 可复用的组件库
- 响应式布局
- 无障碍支持

### 8.2 开发建议

1. **渐进式开发**: 先实现基础功能，再添加视觉效果
2. **性能优化**: 控制动画数量，避免过度渲染
3. **组件复用**: 建立完善的组件库
4. **响应式设计**: 适配不同屏幕尺寸
5. **可访问性**: 保证键盘导航和屏幕阅读器支持

### 8.3 下一步

1. 完善其他页面设计（Deployment列表、Pod详情等）
2. 制作更多业务组件
3. 优化移动端体验
4. 添加深色/浅色主题切换
5. 制作设计规范文档

---

**设计文档版本**: v1.0  
**最后更新**: 2025-11-11  
**设计师**: Claude + Jason
