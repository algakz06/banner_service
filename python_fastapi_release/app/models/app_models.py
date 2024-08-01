from dataclasses import dataclass
from enum import Enum
from typing import Any, Dict, List
import datetime
from uuid import UUID


class UserRoles(Enum):
    admin = "admin"
    user = "user"


@dataclass
class Banner:
    banner_id: int
    is_active: bool
    feature_id: int
    tag_ids: List[int]
    content: Dict[str, Any]
    created_at: datetime.datetime
    updated_at: datetime.datetime


@dataclass
class User:
    id: UUID
    name: str
    password: str
    role: UserRoles
