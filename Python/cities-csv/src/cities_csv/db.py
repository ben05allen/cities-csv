from dotenv import loadenv
import os
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker


loadenv()
SERVER = os.environ.get("DB_SERVER")
DB = os.environ.get("DB_NAME")


def get_session():
    cxn_str = "some db connection string"
    engine = create_engine(cxn_str)

    return sessionmaker(bind=engine)
