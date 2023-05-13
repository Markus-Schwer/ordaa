FROM python:3.9

ADD bot.py .

RUN pip3 install configparser matrix-nio

CMD ["python3", "-u", "./bot.py"]
 
