"""Daily scheduler: retrain + predict at 15:10 (Asia/Shanghai), weekdays only."""
import logging
from datetime import datetime

from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.cron import CronTrigger

from . import config

logger = logging.getLogger("scheduler")


def start_scheduler(service):
    tz = config.settings.timezone
    h, m = config.settings.predict_hour, config.settings.predict_minute
    sched = BackgroundScheduler(timezone=tz)
    sched.add_job(
        service.run_daily,
        CronTrigger(hour=h, minute=m, day_of_week="mon-fri", timezone=tz),
        id="daily_prediction",
        replace_existing=True,
    )
    sched.start()
    logger.info("scheduler started: daily prediction at %02d:%02d %s (weekdays)", h, m, tz)
    return sched
