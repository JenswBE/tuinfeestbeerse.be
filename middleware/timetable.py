from melthon.middleware import Middleware
from collections import OrderedDict
from datetime import date, timedelta, datetime


class Timetable(Middleware):
    def before(self, context):
        # Build time table
        timetable_conf = context['data']['timetable']

        if timetable_conf['settings']['show_timetable']:
            days = []

            # Build lookup table
            for day in timetable_conf['timetable']:
                if day.get('locations'):
                    # Create basic day dict
                    day_dict = {
                        'name': day['name'],
                        'start': datetime.strptime(day['start'], '%d.%m.%Y %H:%M'),
                        'end': datetime.strptime(day['end'], '%d.%m.%Y %H:%M'),
                        'locations': [],
                        'shows': [],
                    }

                    # Fill locations and shows lists
                    for location in day['locations']:
                        day_dict['locations'].append(location['name'])
                        for k, show in enumerate(location['shows']):
                            date_from = (day_dict['end'] if show['from'].startswith('0') else day_dict['start']).strftime('%d.%m.%Y ')
                            date_to = (day_dict['end'] if show['to'].startswith('0') else day_dict['start']).strftime('%d.%m.%Y ')
                            show_dict = {
                                'start': datetime.strptime(date_from + show['from'], '%d.%m.%Y %H:%M'),
                                'end': datetime.strptime(date_to + show['to'], '%d.%m.%Y %H:%M'),
                                'artist': show['artist'],
                                'location': location['name'],
                            }
                            day_dict['shows'].append(show_dict)

                    # Append day to lookup table
                    days.append(day_dict)

            # Build time table
            timetable = OrderedDict()
            half_hour = timedelta(minutes = 30)
            for day in days:
                current_start = day['start']
                slots = []
                while current_start < day['end']:
                    shows = []
                    current_end = current_start + half_hour
                    for location in day['locations']:
                        result = [s for s in day['shows'] if s['location'] == location and s['start'] <= current_start and s['end'] > current_start]
                        if len(result) == 1:
                            show = result[0]
                            if show['start'] == current_start:
                                shows.append({
                                    'artist': show['artist'],
                                    'length': (show['end'] - show['start']) / half_hour,
                                })
                            elif len(result) == 0:
                                shows.append({})
                            else:
                                exit_with_error("Overlapping time slices in shows: {}. Exiting ...".format(result))
                    slots.append({
                        'time': "{} - {}".format(current_start.strftime('%H:%M'), current_end.strftime('%H:%M')),
                        'shows': shows
                    })
                    current_start += half_hour
                timetable[day['name']] = {
                    'locations': day['locations'],
                    'slots': slots,
                }

            context['timetable'] = timetable
        return context
