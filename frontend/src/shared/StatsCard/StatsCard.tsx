import React from 'react';
import './StatsCard.css';

type Color = 'cyan' | 'purple' | 'green';

interface Props {
  icon: React.ReactNode;
  value: number | string;
  label: string;
  status?: string;
  color?: Color;
}

export default function StatsCard({ icon, value, label, status, color = 'cyan' }: Props) {
  return (
    <div className={`stats-card stats-card--${color}`}>
      <div className="stats-card__flow" />
      <div className="stats-card__scan-line" />
      <div className="stats-card__content">
        <div className="stats-card__icon">
          {icon}
          <div className="icon-glow"></div>
        </div>
        <div className="stats-card__value">{value}</div>
        <div className="stats-card__label">{label}</div>
        {status && (
          <div className="stats-card__status">
            <span className="status-dot"></span>
            <span className="status-text">{status}</span>
          </div>
        )}
      </div>
      <div className="corner-decoration corner-tl"></div>
      <div className="corner-decoration corner-tr"></div>
      <div className="corner-decoration corner-bl"></div>
      <div className="corner-decoration corner-br"></div>
    </div>
  );
}

