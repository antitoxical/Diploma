import React from "react";
import { FiHome, FiUser } from "react-icons/fi";
import { GiStack } from 'react-icons/gi';
import { AiOutlineShoppingCart } from 'react-icons/ai';
import { BiHeart } from 'react-icons/bi';
import { IconType } from "react-icons";
import { Link, useLocation } from "react-router-dom";
import "../style/bottomNav.css";

type NavItem = {
  path: string;
  label: string;
  icon: IconType;
};

const navItems: NavItem[] = [
  { path: "/", label: "", icon: FiHome },
  { path: "/catalog", label: "", icon: GiStack },
  { path: "/basket", label: "", icon: AiOutlineShoppingCart },
  { path: "/favorites", label: "", icon: BiHeart },
  { path: "/profile", label: "", icon: FiUser },
];

const IconWrapper = ({ icon: Icon, size = 22 }: { icon: IconType; size?: number }): React.ReactElement | null => {
  return Icon ? Icon({ size }) as React.ReactElement : null;
};

const BottomNav: React.FC = () => {
  const location = useLocation();

  return (
    <nav className="bottom-nav">
      {navItems.map(({ path, label, icon }) => {
        const isActive = location.pathname === path;

        return (
          <Link
            key={path}
            to={path}
            className={`nav-item ${isActive ? "active" : ""}`}
          >
            <IconWrapper icon={icon} size={22} />
            <span>{label}</span>
          </Link>
        );
      })}
    </nav>
  );
};

export default BottomNav;
