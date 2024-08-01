import abc
from typing import List, Optional, Union

from app.models import app_models


class BannerUseCase(abc.ABC):
    @abc.abstractmethod
    def get_user_banner(
        self, tag_ids: List[int], feature_id: int
    ) -> Union[app_models.Banner, None]:
        pass

    @abc.abstractmethod
    def get_banners(
        self,
        tag_id: int,
        feature_id: int,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
    ) -> List[app_models.Banner]:
        pass

    @abc.abstractmethod
    def create_banner(
        self,
        banner: app_models.Banner,
        user: app_models.User,
    ) -> None:
        pass

    @abc.abstractmethod
    def delete_banner(self, banner_id: int) -> None:
        pass

    @abc.abstractmethod
    def update_banner(self, banner: app_models.Banner) -> None:
        pass
