#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdarg.h>

int user_main(void);
int custom_scanf(const char *str, ...);
void read_input_file(const char *filename);
void cleanup(void);

int custom_scanf(const char *str, ...)
{
    char token[100];
    int k = 0;

    va_list ptr;
    va_start(ptr, str);

    for (int i = 0; str[i] != '\0'; i++)
    {
        token[k++] = str[i];
        if (str[i + 1] == '%' || str[i + 1] == '\0')
        {
            token[k] = '\0';
            k = 0;
            char ch1 = token[1];
            if (ch1 == 'i' || ch1 == 'd' || ch1 == 'u')
            {
                int *value = va_arg(ptr, int *);
                fscanf(stdin, "%i", value);
                printf("%d ", *value);
            }
            else if (ch1 == 'h')
            {
                short *value = va_arg(ptr, short *);
                fscanf(stdin, "%hi", value);
                printf("%hi ", *value);
            }
            else if (ch1 == 'c')
            {
                char c; 
                while ((c = fgetc(stdin)) == '\n'
                       || c == ' ' || c == EOF) { 
                } 
                *va_arg(ptr, char*) = c; 
                printf("%c ", c);
            }
            else if (ch1 == 'f')
            {
                float *value = va_arg(ptr, float *);
                fscanf(stdin, "%f", value);
                printf("%f ", *value);
            }

            else if (ch1 == 'l')
            {
                char ch2 = token[2];
                if (ch2 == 'u' || ch2 == 'd' || ch2 == 'i')
                {
                    long *value = va_arg(ptr, long *);
                    fscanf(stdin, "%li", value);
                    printf("%li ", *value);
                }
                else if (ch2 == 'f')
                {
                    double *value = va_arg(ptr, double *);
                    fscanf(stdin, "%lf", value);
                    printf("%lf ", *value);
                }
            }

            else if (ch1 == 'L')
            {
                char ch2 = token[2];
                if (ch2 == 'u' || ch2 == 'd' || ch2 == 'i')
                {
                    long long *value = va_arg(ptr, long long *);
                    fscanf(stdin, "%Li", value);
                    printf("%Ld ", *value);
                }
                else if (ch2 == 'f')
                {
                    long double *value = va_arg(ptr, long double *);
                    fscanf(stdin, "%Lf", value);
                    printf("%Lf ", *value);
                }
            }
            else if (ch1 == 's')
            {
                char *value = va_arg(ptr, char *);
                fscanf(stdin, "%s", value);
                printf("%s ", value);
            }
        }
    }
    printf("\n");
    va_end(ptr);
    return 0;
}

void read_input_file(const char *filename)
{
    freopen(filename, "r", stdin);
    if (!stdin)
    {
        fprintf(stderr, "Error: Cannot open input file %s\n", filename);
        exit(1);
    }
}

void cleanup(void)
{
    fclose(stdin);
}

int main(int argc, char *argv[])
{
    if (argc != 2)
    {
        fprintf(stderr, "Usage: %s <input_file>\n", argv[0]);
        return 1;
    }

    read_input_file(argv[1]);
    int result = user_main();
    cleanup();
    return result;
}

#define scanf custom_scanf