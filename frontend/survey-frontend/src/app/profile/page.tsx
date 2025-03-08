'use client';

import React, { useState, useEffect } from 'react';
import { fetchUserProfile, saveUserProfile } from '@/api-client/profile';
import Input from '@/components/Input';
import Select from '@/components/Select';
import Button from '@/components/Button';
import { getUser } from '@/api-client/auth';

const ProfilePage = () => {
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    birthDate: '',
    phoneNumber: '',
    language: 'en',
  });

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const profile = await getUser()
        setFormData(profile);
      } catch (error) {
        console.error('Failed to load profile:', error);
      }
    };

    loadProfile();
  }, []);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await saveUserProfile(formData);
      alert('Profile saved successfully!');
    } catch (error) {
      console.error('Failed to save profile:', error);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow-md w-full max-w-lg">
        <h1 className="text-xl font-bold mb-4">Edit Profile</h1>
        <Input label="First Name" type="text" name="firstName" value={formData.firstName} onChange={handleInputChange} />
        <Input label="Last Name" type="text" name="lastName" value={formData.lastName} onChange={handleInputChange} />
        <Input label="Birth Date" type="date" name="birthDate" value={formData.birthDate} onChange={handleInputChange} />
        <Input label="Phone Number" type="tel" name="phoneNumber" value={formData.phoneNumber} onChange={handleInputChange} />
        <Select
          label="Language"
          name="language"
          value={formData.language}
          options={[
            { value: 'en', label: 'English' },
            { value: 'ru', label: 'Русский' },
          ]}
          onChange={handleInputChange}
        />
        <Button>Save</Button>
      </form>
    </div>
  );
};

export default ProfilePage;