from melthon.middleware import Middleware
from collections import OrderedDict
from datetime import date, timedelta, datetime


class Timetable(Middleware):
    def before(self, context):
        # Build time table
        timetable_conf = context['data']['timetable']
        if not timetable_conf['settings']['show_timetable']:
            return context

        # Build lookup table
        days = []
        for day in timetable_conf['timetable']:
            if day.get('locations'):
                # Create basic day dict
                day_dict = {
                    'name': day['name'],
                    'start': datetime.strptime(day['start'], '%d.%m.%Y %H:%M'),
                    'end': datetime.strptime(day['end'], '%d.%m.%Y %H:%M'),
                    'slots': [],
                    'locations': [],
                }

                # Generate slots
                current_slot = day_dict['start']
                while current_slot <= day_dict['end']:
                    day_dict['slots'].append(current_slot.strftime("%H:%M"))
                    current_slot += timedelta(minutes=30)

                # Fill locations and shows lists
                for location in day['locations']:
                    shows = []
                    for k, show in enumerate(location['shows']):
                        date_from = (day_dict['end'] if show['from'].startswith(
                            '0') else day_dict['start']).strftime('%d.%m.%Y ')
                        date_to = (day_dict['end'] if show['to'].startswith(
                            '0') else day_dict['start']).strftime('%d.%m.%Y ')
                        show_dict = {
                            'start': datetime.strptime(date_from + show['from'], '%d.%m.%Y %H:%M'),
                            'end': datetime.strptime(date_to + show['to'], '%d.%m.%Y %H:%M'),
                            'artist': show['artist'],
                            'location': location['name'],
                        }
                        show_dict['start_pixels'] = int((
                            show_dict['start'] - day_dict['start']).total_seconds() / 60.0)
                        show_dict['height_pixels'] = int((
                            show_dict['end'] - show_dict['start']).total_seconds() / 60.0)
                        shows.append(show_dict)
                    day_dict['locations'].append({
                        'name': location['name'],
                        'shows': shows,
                    })

                # Append day to lookup table
                days.append(day_dict)
        context['timetable'] = days
        return context
