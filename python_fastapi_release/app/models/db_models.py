import datetime
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy import Column, Integer, String, DateTime, Boolean, ForeignKey, JSON


class Base(DeclarativeBase):
    pass


class AccessLevel(Base):
    __tablename__ = "access_level"
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False)


class Banner(Base):
    __tablename__ = "banner"
    id = Column(Integer, primary_key=True, autoincrement=True)
    feature_id = Column(Integer, nullable=False)
    content = Column(JSON, nullable=False)
    version_number = Column(Integer, default=1)
    original_banner_id = Column(Integer, default=None)
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime, default=datetime.datetime.now())


class User(Base):
    __tablename__ = "user"
    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String)
    hashed_passwrod = Column(String)
    role = Column(Integer, ForeignKey("access_level.id"), default=1)
